package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/params"
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
		log.Warn().Err(err).Msg("could not process message")
		return err
	}
	var message string
	if err != nil {
		log.Error().Err(err).Msg("could not process message")
		message = err.Error()
		err = p.failure(m.Body, message)
	}
	if err != nil {
		log.Fatal().Err(err).Str("message", message).Msg("could not persist failure")
		return err
	}

	log.Trace().Msg("message processed")

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
		return fmt.Errorf("could not execute lambda: %w", err)
	}

	// Then we can decode the result, which contains all the data we need to process.
	var result results.Parsing
	err = json.Unmarshal(output.Payload, &result)
	if err != nil {
		return fmt.Errorf("could not decode parsing result: %w", err)
	}

	// As the Lambda has no access to the DB, we need to fill in the collection ID
	// for all addition and modification jobs here. This could be done in a batch
	// request to improve performance, but it shouldn't have a big impact.
	for _, addition := range result.Additions {
		collection, err := p.collections.One(addition.ChainID, addition.ContractAddress)
		if err != nil {
			return fmt.Errorf("could not get collection for addition: %w", err)
		}
		addition.CollectionID = collection.ID
	}
	for _, modification := range result.Modifications {
		collection, err := p.collections.One(modification.ChainID, modification.ContractAddress)
		if err != nil {
			return fmt.Errorf("could not get collection for modification: %w", err)
		}
		modification.CollectionID = collection.ID
	}

	// As a first step of processing the results, we want to push all of the addition
	// jobs into the addition pipeline, so that the addition dispatcher can start
	// its work as soon as possible.
	payloads := make([][]byte, 0, len(result.Additions))
	for _, addition := range result.Additions {
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

	// In a second step, we insert all of the transfers and sales that were parsed
	// into the events database.
	err = p.transfers.Upsert(result.Transfers...)
	if err != nil {
		return fmt.Errorf("could not upsert transfers: %w", err)
	}
	err = p.sales.Upsert(result.Sales...)
	if err != nil {
		return fmt.Errorf("could not upsert sales: %w", err)
	}

	// Finally, we make sure that we have all of the NFTs that were modified in the
	// database already, at least as placeholders, and we apply the ownership changes
	// for them in one batch operation.
	err = p.nfts.Touch(result.Modifications...)
	if err != nil {
		return fmt.Errorf("could not touch NFTs: %w", err)
	}
	err = p.owners.Change(result.Modifications...)
	if err != nil {
		return fmt.Errorf("could not change owners: %w", err)
	}

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
