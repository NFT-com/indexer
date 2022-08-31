package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/models/failures"
	"github.com/NFT-com/indexer/models/graph"
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
	cfg         AdditionConfig
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
	options ...AdditionOption,
) *AdditionStage {

	cfg := DefaultAdditionConfig
	for _, option := range options {
		option(&cfg)
	}

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
		cfg:         cfg,
	}

	return &a
}

func (a *AdditionStage) HandleMessage(msg *nsq.Message) error {

	// We process the job and if there was no error at all, we just return `nil`, so
	// NSQ marks the message as successfully processed.
	err := a.process(msg.Body)
	if err == nil {
		return nil
	}

	// In most other code paths, we need to decode the addition job for further
	// processing and/or logging, so we just do it once here.
	var addition jobs.Addition
	err = json.Unmarshal(msg.Body, &addition)
	if err != nil {
		a.log.Fatal().
			Hex("msg_id", msg.ID[:]).
			Int64("msg_timestamp", msg.Timestamp).
			Uint16("msg_attempts", msg.Attempts).
			Str("msg_body", string(msg.Body)).
			Err(err).
			Msg("could not decode addition job")
		return err
	}

	log := a.log.With().
		Str("addition_id", addition.ID).
		Uint64("chain_id", addition.ChainID).
		Str("collection_id", addition.CollectionID).
		Str("contract_address", addition.ContractAddress).
		Str("token_id", addition.TokenID).
		Str("token_standard", addition.TokenStandard).
		Logger()

	// If the Ethereum node API response is too large, we split the job into smaller
	// jobs that we add into the pipeline again. This might delay the processing of
	// these heights, but that's OK for the overall speed-up we will get.
	if failures.TooLarge(err) {
		log.Debug().
			Msg("API response too large, splitting job")
		return a.split(&addition)
	}

	// If we have a temporary error and we have not reached the maximum number of
	// attempts on it yet, we will return an error, which will cause NSQ to requeue
	// the job.
	if !results.Permanent(err) && msg.Attempts < uint16(a.cfg.MaxAttempts) {
		log.Warn().
			Uint16("msg_attempts", msg.Attempts).
			Uint("max_attempts", a.cfg.MaxAttempts).
			Err(err).Msg("could not process job, retrying job")
		return err
	}

	// Finally, if we don't have a temporary error, it's either a permanent one, or
	// we have reached the maximum number of attempts. We handle both of these the
	// same, but we log a different message.
	if results.Permanent(err) {
		log.Error().
			Err(err).
			Msg("permanent error encountered, discarding job")
	} else {
		log.Error().
			Uint16("attempts", msg.Attempts).
			Uint("maximum", a.cfg.MaxAttempts).
			Err(err).
			Msg("maximum number of attempts reached, discarding job")
	}

	// In any case, if we don't want to execute the job again, either because it's
	// permanent, or we reached the maximum number of attempts, we try to store it
	// in the database for analytics / review purposes. If that fails, we just crash
	// the service.
	message := err.Error()
	err = a.fail(&addition, message)
	if err != nil {
		a.log.Fatal().Err(err).Str("message", message).Msg("could not persist addition failure")
		return err
	}

	return nil
}

func (a *AdditionStage) process(payload []byte) error {

	// If we are doing a dry run, we skip any kind of processing.
	if a.cfg.DryRun {
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
	if execErr != nil && strings.Contains(execErr.Error(), "URI query for nonexistent token") {
		return a.delete(payload)
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
	err = a.nfts.Upsert(result.NFT)
	if err != nil {
		return fmt.Errorf("could not insert NFT: %w", err)
	}
	err = a.traits.Upsert(result.Traits...)
	if err != nil {
		return fmt.Errorf("could not insert traits: %w", err)
	}

	a.log.Info().
		Str("job_id", result.Job.ID).
		Uint64("chain_id", result.Job.ChainID).
		Str("contract_address", result.Job.ContractAddress).
		Str("token_id", result.Job.TokenID).
		Str("token_standard", result.Job.TokenStandard).
		Int("traits", len(result.Traits)).
		Msg("addition job processed")

	// As we can't know in advance how many requests a Lambda will make, we will
	// wait here to take as many slots on the rate limiter as were needed.
	for i := 0; i < int(result.Requests); i++ {
		a.limit.Take()
	}

	return nil
}

func (a *AdditionStage) split(addition *jobs.Addition) error {

	return nil
}

func (a *AdditionStage) fail(addition *jobs.Addition, message string) error {
	return a.failures.Addition(addition, message)
}

func (a *AdditionStage) delete(payload []byte) error {

	// Decode the payload into the addition job for the deleted NFT.
	var addition jobs.Addition
	err := json.Unmarshal(payload, &addition)
	if err != nil {
		return fmt.Errorf("could not decode addition job: %w", err)
	}

	// Create the dummy NFT to mark it deleted in the database.
	deletion := graph.NFT{
		ID:           addition.NFTID(),
		CollectionID: addition.CollectionID,
		TokenID:      addition.TokenID,
	}
	err = a.nfts.Delete(&deletion)
	if err != nil {
		return fmt.Errorf("could not delete NFT: %w", err)
	}

	return nil
}
