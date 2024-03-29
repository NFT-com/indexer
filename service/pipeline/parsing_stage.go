package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gammazero/deque"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	cache "github.com/Code-Hex/go-generics-cache"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/failures"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type ParsingStage struct {
	ctx         context.Context
	log         zerolog.Logger
	lambda      *lambda.Client
	name        string
	collections CollectionStore
	transfers   TransferStore
	sales       SaleStore
	nfts        NFTStore
	owners      OwnerStore
	failures    FailureStore
	publisher   Publisher
	cache       *cache.Cache[string, struct{}]
	limit       ratelimit.Limiter
	cfg         ParsingConfig
}

func NewParsingStage(
	ctx context.Context,
	log zerolog.Logger,
	lambda *lambda.Client,
	name string,
	transfers TransferStore,
	sales SaleStore,
	collections CollectionStore,
	nfts NFTStore,
	owners OwnerStore,
	failures FailureStore,
	publisher Publisher,
	limit ratelimit.Limiter,
	options ...ParsingOption,
) *ParsingStage {

	cfg := DefaultParsingConfig
	for _, option := range options {
		option(&cfg)
	}

	p := ParsingStage{
		ctx:         ctx,
		log:         log,
		lambda:      lambda,
		name:        name,
		transfers:   transfers,
		sales:       sales,
		collections: collections,
		nfts:        nfts,
		owners:      owners,
		failures:    failures,
		publisher:   publisher,
		limit:       limit,
		cache:       cache.New[string, struct{}](),
		cfg:         cfg,
	}

	return &p
}

func (p *ParsingStage) HandleMessage(msg *nsq.Message) error {

	// We need the decoded job for processing and logging
	var parsing jobs.Parsing
	err := json.Unmarshal(msg.Body, &parsing)
	if err != nil {
		p.log.Fatal().
			Err(err).
			Hex("msg_id", msg.ID[:]).
			Int64("msg_timestamp", msg.Timestamp).
			Uint16("msg_attempts", msg.Attempts).
			Str("msg_body", string(msg.Body)).
			Msg("could not decode parsing job")
	}

	log := p.log.With().
		Str("parsing_id", parsing.ID).
		Uint64("chain_id", parsing.ChainID).
		Uint64("start_height", parsing.StartHeight).
		Uint64("end_height", parsing.EndHeight).
		Strs("contract_addresses", parsing.ContractAddresses).
		Strs("event_hashes", parsing.EventHashes).
		Logger()

	// We process the job and if there was no error at all, we just return `nil`, so
	// NSQ marks the message as successfully processed.
	err = p.process(msg.Body)
	if err == nil {
		return nil
	}

	// If the Ethereum node API response is too large, we split the job into smaller
	// jobs that we add into the pipeline again. This might delay the processing of
	// these heights, but that's OK for the overall speed-up we will get.
	if failures.TooLarge(err) {
		log.Warn().
			Err(err).
			Msg("API response too large, splitting job")

		// Reduce maximum height and addresses according to what we have too much of.
		heights, addresses := parsing.Heights(), parsing.Addresses()
		switch {
		case heights > 1:
			heights = heights / p.cfg.SplitRatio
		case addresses > 1:
			addresses = addresses / p.cfg.SplitRatio
		default:
			// just requeue
			return err
		}

		// Make sure we didn't round down to zero.
		if heights == 0 {
			heights = 1
		}
		if addresses == 0 {
			addresses = 1
		}

		parsings := parsing.Split(heights, addresses)
		err = p.publish(parsings...)
		if err != nil {
			log.Fatal().Err(err).Msg("could not publish split jobs")
		}

		return nil
	}

	// If we have a permanent error, we store the error in the database and we return
	// `nil` to NSQ to signal that it should not be requeued.
	if failures.Permanent(err) {
		log.Error().
			Err(err).
			Msg("aborting job")

		err = p.failures.Parsing(&parsing, err.Error())
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("could not persist parsing failure")
		}
		return nil
	}

	// If we have reached the maximum number of attempts, we also store the error in
	// the database and return `nil` to NSQ to stop further attempts.
	if msg.Attempts >= uint16(p.cfg.MaxAttempts) {
		log.Error().
			Err(err).
			Uint16("msg_attempts", msg.Attempts).
			Uint16("max_attempts", p.cfg.MaxAttempts).
			Msg("discarding job")

		err = p.failures.Parsing(&parsing, err.Error())
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("could not persist parsing failure")
		}
		return nil
	}

	// At this point, we simply want to retry the job, because we it's not a permanent
	// failure and we have not reached the maximum attempts yet.
	log.Warn().
		Err(err).
		Uint16("msg_attempts", msg.Attempts).
		Uint16("max_attempts", p.cfg.MaxAttempts).
		Msg("retrying job")

	return err
}

func (p *ParsingStage) process(payload []byte) error {

	// If we are doing a dry-run, we skip all of the processing here.
	if p.cfg.DryRun {
		return nil
	}

	// Take one slot from the rate limiter to make sure we don't overload the
	// Ethereum node's API.
	p.limit.Take()

	// Invoke the AWS Lambda function that will retrieve the logs and parse them.
	input := &lambda.InvokeInput{
		FunctionName: aws.String(p.name),
		Payload:      payload,
	}
	output, err := p.lambda.Invoke(p.ctx, input)
	if err != nil {
		return fmt.Errorf("could not invoke lambda: %w", err)
	}

	// Next to invocation errors, we need to handle errors that happened during
	// execution, which are read from a separate field.
	var execErr *results.Error
	if output.FunctionError != nil {
		err = json.Unmarshal(output.Payload, &execErr)
	}
	if err != nil {
		return fmt.Errorf("could not decode execution error: %w", err)
	}
	if execErr != nil {
		return fmt.Errorf("could not execute lambda: %w", execErr)
	}

	// Then we can decode the result, which contains all the data we need to process.
	var result results.Parsing
	err = json.Unmarshal(output.Payload, &result)
	if err != nil {
		return fmt.Errorf("could not decode parsing result: %w", err)
	}

	// We can go through the transfers and process those with zero address as mints.
	var touches []*graph.NFT
	var deletions []*graph.NFT
	var additionPayloads [][]byte
	var owners []*events.Transfer
	for _, transfer := range result.Transfers {

		if transfer.SenderAddress == transfer.ReceiverAddress {
			continue
		}

		// Collect transfers that are not no-ops for owner changes.
		owners = append(owners, transfer)

		// Get the collection ID based on chain ID and collection address, so we can
		// reference it directly for the addition job and the NFT insertion.
		collection, err := p.collections.One(transfer.ChainID, transfer.CollectionAddress)
		if err != nil {
			return fmt.Errorf("could not get collection for transfer: %w", err)
		}

		// Create a placeholder NFT that we will create in the DB.
		if !p.cache.Contains(transfer.NFTID()) {
			p.cache.Set(transfer.NFTID(), struct{}{}, cache.WithExpiration(15*time.Minute))
			touch := graph.NFT{
				ID:           transfer.NFTID(),
				CollectionID: collection.ID,
				TokenID:      transfer.TokenID,
			}
			touches = append(touches, &touch)
		}

		// Create addition jobs for transfers that come from the zero address.
		if transfer.SenderAddress == params.AddressZero {
			addition := jobs.Addition{
				ID:              uuid.NewString(),
				ChainID:         transfer.ChainID,
				CollectionID:    collection.ID,
				ContractAddress: transfer.CollectionAddress,
				TokenID:         transfer.TokenID,
				TokenStandard:   transfer.TokenStandard,
			}
			payload, err := json.Marshal(addition)
			if err != nil {
				return fmt.Errorf("could not encode addition job: %w", err)
			}
			additionPayloads = append(additionPayloads, payload)
		}

		// Create deletion jobs for transfers that go to the zero address, or another
		// known burn address.
		if transfer.ReceiverAddress == params.AddressZero || transfer.ReceiverAddress == params.AddressDead {
			deletion := graph.NFT{
				ID:           transfer.NFTID(),
				CollectionID: collection.ID,
				TokenID:      transfer.TokenID,
			}
			deletions = append(deletions, &deletion)
		}

	}

	// Filter the touches so we only create them for NFTs that are not yet in the DB.
	touches, err = p.nfts.Missing(touches...)
	if err != nil {
		return fmt.Errorf("could not filter touches: %w", err)
	}

	// Touch all the NFTs that have been created, so that we can apply owner changes
	// out of order, before the full NFT information is available from the addition.
	err = p.nfts.Touch(touches...)
	if err != nil {
		return fmt.Errorf("could not execute touches: %w", err)
	}

	// Delete all the NFTs that have been deleted, so that we can apply deletions
	// out of order before the full NFT information is available from the addition.
	err = p.nfts.Delete(deletions...)
	if err != nil {
		return fmt.Errorf("could not execute deletions: %w", err)
	}

	// Next, we can store all the raw events for transfers and sales.
	err = p.transfers.Upsert(result.Transfers...)
	if err != nil {
		return fmt.Errorf("could not upsert transfers: %w", err)
	}
	err = p.sales.Upsert(result.Sales...)
	if err != nil {
		return fmt.Errorf("could not upsert sales: %w", err)
	}

	// Last but not least, we can upsert the owner change updates for each transfer.
	err = p.owners.Upsert(owners...)
	if err != nil {
		return fmt.Errorf("could not upsert owners: %w", err)
	}

	// ... and sanitize the owners table every so often.
	dice := rand.Intn(int(p.cfg.SanitizeInterval))
	if dice == 0 {
		err = p.owners.Sanitize()
		if err != nil {
			return fmt.Errorf("could not sanitize owners: %w", err)
		}
	}

	if len(additionPayloads) > 0 {
		err = p.publisher.MultiPublish(params.TopicAddition, additionPayloads)
		if err != nil {
			return fmt.Errorf("could not publish addition jobs: %w", err)
		}
	}

	// We can go through the sales and process the completion.
	salesMap := make(map[uint64][]*events.Sale)
	for _, sale := range result.Sales {
		if !sale.NeedsCompletion {
			continue
		}
		salesMap[sale.BlockNumber] = append(salesMap[sale.BlockNumber], sale)
	}

	var completionPayloads [][]byte
	for height, sales := range salesMap {
		completion := jobs.Completion{
			ID:          uuid.NewString(),
			ChainID:     result.Job.ChainID,
			StartHeight: height,
			EndHeight:   height,
			Sales:       sales,
		}
		payload, err := json.Marshal(completion)
		if err != nil {
			return fmt.Errorf("could not encode completion job: %w", err)
		}
		completionPayloads = append(completionPayloads, payload)
	}

	if len(completionPayloads) > 0 {
		err = p.publisher.MultiPublish(params.TopicCompletion, completionPayloads)
		if err != nil {
			return fmt.Errorf("could not publish completion job: %w", err)
		}
	}

	p.log.Info().
		Str("job_id", result.Job.ID).
		Uint64("chain_id", result.Job.ChainID).
		Uint64("start_height", result.Job.StartHeight).
		Uint64("end_height", result.Job.EndHeight).
		Strs("contract_addresses", result.Job.ContractAddresses).
		Strs("event_hashes", result.Job.EventHashes).
		Int("transfers", len(result.Transfers)).
		Int("sales", len(result.Sales)).
		Msg("parsing job processed")

	return nil
}

func (p *ParsingStage) publish(parsings ...*jobs.Parsing) error {

	var queue deque.Deque[[]*jobs.Parsing]
	queue.PushBack(parsings)
	for queue.Len() > 0 {

		parsings = queue.PopFront()
		payloads := make([][]byte, 0, len(parsings))
		for _, parsing := range parsings {
			payload, err := json.Marshal(parsing)
			if err != nil {
				return fmt.Errorf("could not encode parsing job: %w", err)
			}
			payloads = append(payloads, payload)
		}

		err := p.publisher.MultiPublish(params.TopicParsing, payloads)
		if err == nil {
			continue
		}

		if strings.Contains(err.Error(), "MPUB body too big") {
			end := len(parsings)
			pivot := end / 2
			part1 := parsings[0:pivot]
			part2 := parsings[pivot:end]
			queue.PushBack(part1)
			queue.PushBack(part2)
			continue
		}

		if err != nil {
			return fmt.Errorf("could not publish parsing jobs: %w", err)
		}
	}

	return nil
}
