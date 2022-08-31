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
	cfg         CompletionConfig
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
	options ...CompletionOption,
) *CompletionStage {

	cfg := DefaultCompletionConfig
	for _, option := range options {
		option(&cfg)
	}

	a := CompletionStage{
		ctx:         ctx,
		log:         log,
		lambda:      lambda,
		name:        name,
		collections: collections,
		sales:       sales,
		failures:    failures,
		limit:       limit,
		cfg:         cfg,
	}

	return &a
}

func (c *CompletionStage) HandleMessage(m *nsq.Message) error {

	err := c.process(m.Body)
	if err == nil {
		return nil
	}

	// We only retry if we don't have a permanent error, and we have not reached
	// the maximum number of attempts.
	if !results.Permanent(err) && m.Attempts < uint16(c.cfg.MaxAttempts) {
		c.log.Warn().
			Uint16("attempts", m.Attempts).
			Uint("maximum", c.cfg.MaxAttempts).
			Err(err).Msg("could not process job, retrying")
		return err
	}

	// Otherwise, we either have a permanent error, or maximum number of retries.
	// Log the error accordingly, and proceed without retrying (return `nil`).
	if results.Permanent(err) {
		c.log.Error().Err(err).Msg("permanent error encountered, aborting")
	} else {
		c.log.Error().
			Uint16("attempts", m.Attempts).
			Uint("maximum", c.cfg.MaxAttempts).
			Err(err).
			Msg("maximum number of attempts reached, aborting")
	}

	// The below code stores the error in the database and fatally fails the service
	// if we don't manage to do so.
	message := err.Error()
	err = c.failure(m.Body, message)
	if err != nil {
		c.log.Fatal().Err(err).Str("message", message).Msg("could not persist completion failure")
		return err
	}

	return nil
}

func (c *CompletionStage) process(payload []byte) error {

	// If we are doing a dry run, we skip any kind of processing.
	if c.cfg.DryRun {
		return nil
	}

	// We then take up a slot in the rate limiter to make sure the Ethereum node
	// API is not overloaded.
	c.limit.Take()

	// Next, we invoke the Lambda, which will get the token URI, query it and parse it.
	input := &lambda.InvokeInput{
		FunctionName: aws.String(c.name),
		Payload:      payload,
	}
	output, err := c.lambda.Invoke(c.ctx, input)
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
		return fmt.Errorf("could not decode completion result: %w", err)
	}

	// Finally, we make the necessary changes to the DB: update sale event.
	err = c.sales.Update(result.Job.Sales...)
	if err != nil {
		return fmt.Errorf("could not update sale event: %w", err)
	}

	c.log.Info().
		Str("job_id", result.Job.ID).
		Uint64("chain_id", result.Job.ChainID).
		Uint64("start_height", result.Job.StartHeight).
		Uint64("end_height", result.Job.EndHeight).
		Int("sales", len(result.Job.Sales)).
		Msg("completion job processed")

	return nil
}

func (c *CompletionStage) failure(payload []byte, message string) error {

	// Decode the payload into the failed completion job.
	var completion jobs.Completion
	err := json.Unmarshal(payload, &completion)
	if err != nil {
		return fmt.Errorf("could not decode completion job: %w", err)
	}

	// Persist the completion failure in the DB, so it can be reviewed and potentially
	// retried at a later point.
	err = c.failures.Completion(&completion, message)
	if err != nil {
		return fmt.Errorf("could not persist completion failure: %w", err)
	}

	return nil
}
