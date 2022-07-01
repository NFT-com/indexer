package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type CompletionStage struct {
	ctx         context.Context
	log         zerolog.Logger
	lambda      *lambda.Client
	name        string
	collections CollectionStore
	sales       SaleStore
	failures    FailureStore
	limit       ratelimit.Limiter
	dryRun      bool
}

func NewCompletionStage(
	ctx context.Context,
	log zerolog.Logger,
	lambda *lambda.Client,
	name string,
	collections CollectionStore,
	sales SaleStore,
	failures FailureStore,
	limit ratelimit.Limiter,
	dryRun bool,
) *CompletionStage {

	a := CompletionStage{
		ctx:         ctx,
		log:         log,
		lambda:      lambda,
		name:        name,
		collections: collections,
		sales:       sales,
		failures:    failures,
		limit:       limit,
		dryRun:      dryRun,
	}

	return &a
}

func (a *CompletionStage) HandleMessage(m *nsq.Message) error {

	err := a.process(m.Body)
	if results.Retriable(err) {
		a.log.Warn().Err(err).Msg("could not process message, retrying")
		return err
	}
	var message string
	if err != nil {
		a.log.Error().Err(err).Msg("could not process message, discarding")
		message = err.Error()
		err = a.failure(m.Body, message)
	}
	if err != nil {
		a.log.Fatal().Err(err).Str("message", message).Msg("could not persist completion failure")
		return err
	}

	return nil
}

func (a *CompletionStage) process(payload []byte) error {

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
	var result results.Completion
	err = json.Unmarshal(output.Payload, &result)
	if err != nil {
		return fmt.Errorf("could not decode NFT: %w", err)
	}

	// Finally, we make the necessary changes to the DB: update sale event.
	err = a.sales.Update(result.Sale)
	if err != nil {
		return fmt.Errorf("could not insert NFT: %w", err)
	}

	a.log.Info().
		Str("job_id", result.Job.ID).
		Uint64("chain_id", result.Job.ChainID).
		Str("sale_id", result.Job.SaleID).
		Str("collection_address", result.Sale.CollectionAddress).
		Str("token_id", result.Sale.CollectionAddress).
		Msg("completion job processed")

	// As we can't know in advance how many requests a Lambda will make, we will
	// wait here to take as many slots on the rate limiter as were needed.
	for i := 0; i < int(result.Requests); i++ {
		a.limit.Take()
	}

	return nil
}

func (a *CompletionStage) failure(payload []byte, message string) error {

	// Decode the payload into the failed completion job.
	var completion jobs.Completion
	err := json.Unmarshal(payload, &completion)
	if err != nil {
		return fmt.Errorf("could not decode completion job: %w", err)
	}

	// Persist the completion failure in the DB, so it can be reviewed and potentially
	// retried at a later point.
	err = a.failures.Completion(&completion, message)
	if err != nil {
		return fmt.Errorf("could not persist completion failure: %w", err)
	}

	return nil
}
