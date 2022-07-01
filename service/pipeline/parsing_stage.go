package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

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
	additions   BatchPublisher
	limit       ratelimit.Limiter
	dryRun      bool
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
	additions BatchPublisher,
	limit ratelimit.Limiter,
	dryRun bool,
) *ParsingStage {

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
		additions:   additions,
		limit:       limit,
		dryRun:      dryRun,
	}

	return &p
}

func (p *ParsingStage) HandleMessage(m *nsq.Message) error {

	err := p.process(m.Body)
	if results.Retriable(err) {
		p.log.Warn().Err(err).Msg("could not process message, retrying")
		return err
	}
	var message string
	if err != nil {
		p.log.Error().Err(err).Msg("could not process message, discarding")
		message = err.Error()
		err = p.failure(m.Body, message)
	}
	if err != nil {
		p.log.Fatal().Err(err).Str("message", message).Msg("could not persist parsing failure")
		return err
	}

	return nil
}

func (p *ParsingStage) process(payload []byte) error {

	// If we are doing a dry-run, we skip all of the processing here.
	if p.dryRun {
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
	var dummies []*graph.NFT
	var payloads [][]byte
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
		dummy := graph.NFT{
			ID:           transfer.NFTID(),
			CollectionID: collection.ID,
			TokenID:      transfer.TokenID,
		}
		dummies = append(dummies, &dummy)

		// Skip transfers that do not originate from the zero address so we process
		// only mints.
		if transfer.SenderAddress != params.AddressZero {
			continue
		}

		// Create an addition job to complete the data for the NFT.
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
		payloads = append(payloads, payload)
	}

	if len(payloads) > 0 {
		err = p.additions.MultiPublish(params.TopicAddition, payloads)
		if err != nil {
			return fmt.Errorf("could not publish addition jobs: %w", err)
		}
	}

	// Touch all the NFTs that have been created, so that we can apply owner changes
	// out of order, before the full NFT information is available from the addition.
	err = p.nfts.Touch(dummies...)
	if err != nil {
		return fmt.Errorf("could not touch dummies: %w", err)
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

	// Persist the parsing failure in the DB so it can be reviewed and potentially
	// retried at a later point.
	err = p.failures.Parsing(&parsing, message)
	if err != nil {
		return fmt.Errorf("could not persist parsing failure: %w", err)
	}

	return nil
}
