package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/models/events"
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
		cfg:         cfg,
	}

	return &p
}

func (p *ParsingStage) HandleMessage(m *nsq.Message) error {

	err := p.process(m.Body)
	if err == nil {
		return nil
	}

	if m.Attempts >= uint16(p.cfg.MaxRetries) {
		p.log.Error().Err(err).Msg("maximum number of retries reached, aborting")
		return err
	}

	if !results.Permanent(err) {
		p.log.Warn().Err(err).Msg("temporary error encountered, retrying")
		return err
	}

	p.log.Error().Err(err).Msg("permanent error encountered, discarding")

	message := err.Error()
	err = p.failure(m.Body, message)
	if err != nil {
		p.log.Fatal().Err(err).Str("message", message).Msg("could not persist parsing failure")
		return err
	}

	return nil
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
		touch := graph.NFT{
			ID:           transfer.NFTID(),
			CollectionID: collection.ID,
			TokenID:      transfer.TokenID,
		}
		touches = append(touches, &touch)

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

	// As we can't know in advance how many requests a Lambda will make, we will
	// wait here to take up any requests above one that we needed.
	for i := 1; i < int(result.Requests); i++ {
		p.limit.Take()
	}

	return nil
}

func (p *ParsingStage) failure(payload []byte, message string) error {

	// Decode the payload into the failed parsing job.
	var parsing jobs.Parsing
	err := json.Unmarshal(payload, &parsing)
	if err != nil {
		return fmt.Errorf("could not decode parsing job: %w", err)
	}

	// Persist the parsing failure in the DB, so it can be reviewed and potentially
	// retried at a later point.
	err = p.failures.Parsing(&parsing, message)
	if err != nil {
		return fmt.Errorf("could not persist parsing failure: %w", err)
	}

	return nil
}
