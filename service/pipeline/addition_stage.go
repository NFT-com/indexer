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

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type AdditionStage struct {
	ctx         context.Context
	log         zerolog.Logger
	lambda      *lambda.Client
	name        string
	collections CollectionStore
	nfts        NFTStore
	owners      OwnerStore
	traits      TraitStore
	failures    FailureStore
	limit       ratelimit.Limiter
	dryRun      bool
}

func NewAdditionStage(
	ctx context.Context,
	log zerolog.Logger,
	lambda *lambda.Client,
	name string,
	collections CollectionStore,
	nfts NFTStore,
	owners OwnerStore,
	traits TraitStore,
	failures FailureStore,
	limit ratelimit.Limiter,
	dryRun bool,
) *AdditionStage {

	a := AdditionStage{
		ctx:         ctx,
		log:         log,
		lambda:      lambda,
		name:        name,
		collections: collections,
		nfts:        nfts,
		owners:      owners,
		traits:      traits,
		failures:    failures,
		limit:       limit,
		dryRun:      dryRun,
	}

	return &a
}

func (a *AdditionStage) HandleMessage(m *nsq.Message) error {

	err := a.process(m.Body)
	if results.Retriable(err) {
		log.Warn().Err(err).Msg("could not process message, retrying")
		return err
	}
	var message string
	if err != nil {
		log.Error().Err(err).Msg("could not process message, discarding")
		message = err.Error()
		err = a.failure(m.Body, message)
	}
	if err != nil {
		log.Fatal().Err(err).Str("message", message).Msg("could not persist addition failure")
		return err
	}

	return nil
}

func (a *AdditionStage) process(payload []byte) error {

	// If we are doing a dry run, we skip any kind of processing.
	if a.dryRun {
		return nil
	}

	// We then take up a slot in the rate limiter to make sure the Ethereum node
	// API is not overloaded.
	a.limit.Take()

	// Next, we invoke the Lambda, which will get the token URI, query it and parse it.
	input := &lambda.InvokeInput{
		FunctionName: aws.String(a.name),
		Payload:      payload,
	}
	output, err := a.lambda.Invoke(a.ctx, input)
	if err != nil {
		return fmt.Errorf("could not invoke lambda: %w", err)
	}

	// Execution errors are handled separately from invocation errors, as they are
	// returned as part of the result.
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

	// We then unmarshal the result.
	var result results.Addition
	err = json.Unmarshal(output.Payload, &result)
	if err != nil {
		return fmt.Errorf("could not decode NFT: %w", err)
	}

	// Finally, we make the necessary changes to the DB: insert the NFT, the traits
	// and apply the necessary ownership changes.
	err = a.nfts.Insert(result.NFT)
	if err != nil {
		return fmt.Errorf("could not insert NFT: %w", err)
	}
	err = a.traits.Insert(result.Traits...)
	if err != nil {
		return fmt.Errorf("could not insert traits: %w", err)
	}
	err = a.owners.Add(&result)
	if err != nil {
		return fmt.Errorf("could not add owner: %w", err)
	}

	a.log.Info().
		Str("job_id", result.Job.ID).
		Str("contract_address", result.Job.ContractAddress).
		Str("token_id", result.Job.TokenID).
		Str("token_standard", result.Job.TokenStandard).
		Str("owner_address", result.Job.OwnerAddress).
		Uint("token_count", result.Job.TokenCount).
		Str("collection_id", result.NFT.CollectionID).
		Str("nft_id", result.NFT.ID).
		Int("traits", len(result.Traits)).
		Msg("addition job processed")

	// As we can't know in advance how many requests a Lambda will make, we will
	// wait here to take as many slots on the rate limiter as were needed.
	for i := 0; i < int(result.Requests); i++ {
		a.limit.Take()
	}

	return nil
}

func (a *AdditionStage) failure(payload []byte, message string) error {

	// Decode the payload into the failed addition job.
	var addition jobs.Addition
	err := json.Unmarshal(payload, &addition)
	if err != nil {
		return fmt.Errorf("could not decode addition job: %w", err)
	}

	// Persist the addition failure in the DB so it can be reviewed and potentially
	// retried at a later point.
	err = a.failures.Addition(&addition, message)
	if err != nil {
		return fmt.Errorf("could not persist addition failure: %w", err)
	}

	return nil
}
