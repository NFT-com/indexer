package pipeline

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/function/handlers/parsing"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/logs"
)

type ParsingConsumer struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	parsings      ParsingStore
	mints         MintStore
	transfers     TransferStore
	sales         SaleStore
	burns         BurnStore
	limit         ratelimit.Limiter
	consumerQueue chan [][]byte
	close         chan struct{}
	dryRun        bool
}

func NewParsingConsumer(
	log zerolog.Logger,
	dispatcher function.Invoker,
	parsings ParsingStore,
	mints MintStore,
	transfers TransferStore,
	sales SaleStore,
	burns BurnStore,
	rateLimit int,
	dryRun bool,
) *ParsingConsumer {

	p := ParsingConsumer{
		log:           log,
		dispatcher:    dispatcher,
		parsings:      parsings,
		mints:         mints,
		transfers:     transfers,
		sales:         sales,
		burns:         burns,
		limit:         ratelimit.New(rateLimit),
		consumerQueue: make(chan [][]byte, concurrentConsumers),
		close:         make(chan struct{}),
		dryRun:        dryRun,
	}

	return &p
}

func (p *ParsingConsumer) Consume(batch rmq.Deliveries) {

	p.log.Debug().Int("jobs", len(batch)).Msg("received batch for consuming")

	payloads := make([][]byte, 0, len(batch))

	for _, delivery := range batch {
		payload := []byte(delivery.Payload())
		payloads = append(payloads, payload)

		err := delivery.Ack()
		if err != nil {
			p.log.Error().Err(err).Msg("could not acknowledge message")
			return
		}
	}

	p.consumerQueue <- payloads
}

func (p *ParsingConsumer) Run() {
	for i := 0; i < concurrentConsumers; i++ {
		go func() {
			for {
				select {
				case <-p.close:
					return
				case payload := <-p.consumerQueue:
					p.consume(payload)
				}
			}
		}()
	}
}

func (p *ParsingConsumer) Close() {
	close(p.close)
}

func (p *ParsingConsumer) consume(payloads [][]byte) {

	parsings := p.unmarshalParsings(payloads)

	inputMap := make(map[uint64]jobs.Input, len(parsings))

	// Job List is unordered, in order to group them in a continuous way,
	// first the next loop basically maps the height to an input that could group jobs
	// from the same height. It also gets the highest and lowest heights in the list.
	lowestBlock := uint64(math.MaxUint64)
	highestBlock := uint64(0)
	for _, parsing := range parsings {
		block := parsing.BlockNumber
		// checks if there is already an entry, if so joins the ids, addresses, event_types and maps the standard.
		input, ok := inputMap[block]
		if ok {
			input.IDs = append(input.IDs, parsing.ID)
			input.Addresses = append(input.Addresses, parsing.Address)
			input.EventTypes = append(input.EventTypes, parsing.Event)
			input.Standards[parsing.Event] = parsing.Standard

			inputMap[block] = input
			continue
		}

		// if there is no entry just create a new one and check if this is lower that
		// the current lowest height or upper than the highest height.
		inputMap[block] = jobs.Input{
			IDs:        []string{parsing.ID},
			ChainURL:   parsing.ChainURL,
			ChainID:    parsing.ChainID,
			ChainType:  parsing.ChainType,
			StartBlock: parsing.BlockNumber,
			EndBlock:   parsing.BlockNumber,
			Addresses:  []string{parsing.Address},
			Standards:  map[string]string{parsing.Event: parsing.Standard},
			EventTypes: []string{parsing.Event},
		}

		if lowestBlock > block {
			lowestBlock = block
		}

		if highestBlock < block {
			highestBlock = block
		}
	}

	inputs := make([]jobs.Input, 0)

	// Merges all the continuous inputs heights into a single one.
	input := jobs.Input{}
	for i := lowestBlock; i <= highestBlock; i++ {
		part, ok := inputMap[i]
		if !ok {
			if len(input.IDs) != 0 {
				inputs = append(inputs, input)
			}

			input = jobs.Input{}
			continue
		}

		// If current input ids len is 0, it means it does not have a current input
		if len(input.IDs) == 0 {
			input = part
			continue
		}

		input = p.fillInput(input, part)
	}

	if len(input.IDs) != 0 {
		inputs = append(inputs, input)
	}

	p.log.Debug().Int("jobs", len(payloads)).Int("batches", len(inputs)).Msg("batched jobs for processing")

	for _, input := range inputs {

		payload, err := json.Marshal(input)
		if err != nil {
			p.handleError(err, "could not marshal input payload", input.IDs...)
			return
		}

		err = p.parsings.Update(input.IDs, jobs.StatusProcessing)
		if err != nil {
			p.handleError(err, "could not update jobs statuses")
			return
		}

		// Wait for rate limiter to have available spots.
		p.limit.Take()

		p.log.Debug().
			Uint64("start", input.StartBlock).
			Uint64("end", input.EndBlock).
			Int("collections", len(input.Addresses)).
			Int("standards", len(input.Standards)).
			Int("events", len(input.EventTypes)).
			Msg("dispatching job batch")

		name := functionName(input)

		if !p.dryRun {
			notify := func(err error, dur time.Duration) {
				p.log.Error().
					Err(err).
					Dur("retry_in", dur).
					Str("name", name).
					Int("payload_len", len(payload)).
					Msg("could not invoke lambda")
			}
			var output []byte
			_ = backoff.RetryNotify(func() error {

				output, err = p.dispatcher.Invoke(name, payload)
				if err != nil {
					return err
				}
				return nil
			}, backoff.NewExponentialBackOff(), notify)

			var entries []logs.Entry
			err = json.Unmarshal(output, &entries)
			if err != nil {
				p.handleError(err, "could not unmarshal output log entries", input.IDs...)
				return
			}

			p.log.Debug().
				Uint64("start", input.StartBlock).
				Uint64("end", input.EndBlock).
				Int("collections", len(input.Addresses)).
				Int("standards", len(input.Standards)).
				Int("events", len(input.EventTypes)).
				Int("entries", len(entries)).
				Msg("processing results")

			err = p.processEntries(input, entries)
			if err != nil {
				p.handleError(err, "could not handle output logs", input.IDs...)
				return
			}
		}

		err = p.dataStore.UpdateParsingJobsStatus(input.IDs, jobs.StatusFinished)
		if err != nil {
			p.handleError(err, "could not update jobs statuses")
			return
		}
	}
}

func (p *ParsingConsumer) fillInput(base, part jobs.Input) jobs.Input {
	base.EndBlock = part.EndBlock
	base.IDs = append(base.IDs, part.IDs...)

	addresses := make(map[string]struct{}, len(base.Addresses))
	for _, address := range base.Addresses {
		addresses[address] = struct{}{}
	}

	for _, address := range part.Addresses {
		if _, ok := addresses[address]; ok {
			continue
		}

		base.Addresses = append(base.Addresses, address)
	}

	eventTypes := make(map[string]struct{}, len(base.EventTypes))
	for _, eventType := range base.EventTypes {
		eventTypes[eventType] = struct{}{}
	}

	for _, eventType := range part.EventTypes {
		if _, ok := eventTypes[eventType]; ok {
			continue
		}

		base.EventTypes = append(base.EventTypes, eventType)
	}

	return base
}

func (p *ParsingConsumer) unmarshalParsings(payloads [][]byte) []*jobs.Parsing {
	parsings := make([]*jobs.Parsing, 0, len(payloads))

	for _, payload := range payloads {
		var job jobs.Parsing
		err := json.Unmarshal(payload, &job)
		if err != nil {
			p.log.Error().Err(err).Msg("could not unmarshal message")
			continue
		}

		parsings = append(parsings, &job)
	}

	return parsings
}

func (p *ParsingConsumer) handleError(err error, message string, ids ...string) {
	updateErr := p.jobStore.UpdateParsingJobsStatus(ids, jobs.StatusFailed)
	if updateErr != nil {
		p.log.Error().Err(updateErr).Msg("could not update jobs statuses")
	}

	p.log.Error().Err(err).Strs("job_ids", ids).Msg(message)
}

func blockToUint64(block string) (uint64, error) {
	number := big.NewInt(0)
	_, ok := number.SetString(block, 0)
	if !ok {
		return 0, fmt.Errorf("could not parse block into big.Int")
	}

	return number.Uint64(), nil
}

func functionName(input parsing.Input) string {
	h := sha256.New()

	s := strings.Join(
		[]string{
			"parsing",
			strings.ToLower(input.ChainType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}

func (p *ParsingConsumer) processEntries(input jobs.Input, entries []logs.Entry) error {
	for _, entry := range entries {

		chain, err := p.dataStore.Chain(entry.ChainID)
		if err != nil {
			return fmt.Errorf("could not get chain: %w", err)
		}

		if entry.NeedsActionJob {
			err := p.jobStore.CreateActionJob(&jobs.Action{
				ID:          uuid.New().String(),
				ChainURL:    input.ChainURL,
				ChainID:     input.ChainID,
				ChainType:   input.ChainType,
				BlockNumber: entry.Block,
				Address:     entry.Contract,
				Standard:    entry.Standard,
				Event:       entry.Event,
				TokenID:     entry.NftID,
				ToAddress:   entry.ToAddress,
				Type:        entry.ActionType.String(),
				Status:      jobs.StatusCreated,
			})
			if err != nil {
				return fmt.Errorf("could not create action job: %w", err)
			}
		}

		switch entry.Type {

		case logs.TypeMint:

			collection, err := p.dataStore.Collection(chain.ID, entry.Contract, entry.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection (chainID: %s contract: %s): %w", chain.ChainID, entry.Contract, err)
			}

			event := events.Mint{
				ID:              entry.ID,
				CollectionID:    collection.ID,
				Block:           entry.Block,
				EventIndex:      entry.Index,
				TransactionHash: entry.TransactionHash,
				TokenID:         entry.NftID,
				Owner:           entry.ToAddress,
				EmittedAt:       entry.EmittedAt,
			}

			err = p.eventStore.UpsertMintEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert mint event: %w", err)
			}

		case logs.TypeSale:

			marketplace, err := p.dataStore.Marketplace(chain.ID, entry.Contract)
			if err != nil {
				return fmt.Errorf("could not get marketplace: %w", err)
			}

			event := events.Sale{
				ID:              entry.ID,
				MarketplaceID:   marketplace.ID,
				Block:           entry.Block,
				EventIndex:      entry.Index,
				TransactionHash: entry.TransactionHash,
				Seller:          entry.ToAddress,
				Buyer:           entry.FromAddress,
				Price:           entry.Price,
				EmittedAt:       entry.EmittedAt,
			}

			err = p.eventStore.UpsertSaleEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert sale event: %w", err)
			}

		case logs.TypeTransfer:

			collection, err := p.dataStore.Collection(chain.ID, entry.Contract, entry.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection (chainID: %s contract: %s): %w", chain.ChainID, entry.Contract, err)
			}

			event := events.Transfer{
				ID:              entry.ID,
				CollectionID:    collection.ID,
				Block:           entry.Block,
				EventIndex:      entry.Index,
				TransactionHash: entry.TransactionHash,
				TokenID:         entry.NftID,
				FromAddress:     entry.FromAddress,
				ToAddress:       entry.ToAddress,
				EmittedAt:       entry.EmittedAt,
			}

			err = p.eventStore.UpsertTransferEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert transfer event: %w", err)
			}

		case logs.TypeBurn:

			collection, err := p.dataStore.Collection(chain.ID, entry.Contract, entry.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection (chainID: %s contract: %s): %w", chain.ChainID, entry.Contract, err)
			}

			event := events.Burn{
				ID:              entry.ID,
				CollectionID:    collection.ID,
				Block:           entry.Block,
				EventIndex:      entry.Index,
				TransactionHash: entry.TransactionHash,
				TokenID:         entry.NftID,
				EmittedAt:       entry.EmittedAt,
			}

			err = p.eventStore.UpsertBurnEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert burn event: %w", err)
			}

		default:
			p.log.Error().Str("type", entry.Type.String()).Msg("got unexpected log type")
		}
	}

	return nil
}

type EventStore interface {
	UpsertMintEvent(event events.Mint) error
	UpsertSaleEvent(event events.Sale) error
	UpsertTransferEvent(event events.Transfer) error
	UpsertBurnEvent(event events.Burn) error

	Chain(chainID string) (*graph.Chain, error)
	Collection(chainID, address, contractCollectionID string) (*graph.Collection, error)
	Marketplace(chainID, address string) (*graph.Marketplace, error)

	UpdateNFTOwner(collectionID, nft, owner string) error

	CreateActionJob(job *jobs.Action) error
	ParsingJob(id string) (*jobs.Parsing, error)
	UpdateParsingJobsStatus(ids []string, status string) error
}
