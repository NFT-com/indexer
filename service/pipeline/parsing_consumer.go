package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type ParsingConsumer struct {
	log       zerolog.Logger
	lambda    *lambda.Lambda
	parsings  ParsingStore
	actions   ActionStore
	mints     MintStore
	transfers TransferStore
	sales     SaleStore
	burns     BurnStore
	limit     ratelimit.Limiter
	dryRun    bool
}

func NewParsingConsumer(
	log zerolog.Logger,
	lambda *lambda.Lambda,
	parsings ParsingStore,
	actions ActionStore,
	mints MintStore,
	transfers TransferStore,
	sales SaleStore,
	burns BurnStore,
	rateLimit int,
	dryRun bool,
) *ParsingConsumer {

	p := ParsingConsumer{
		log:       log,
		lambda:    lambda,
		parsings:  parsings,
		actions:   actions,
		mints:     mints,
		transfers: transfers,
		sales:     sales,
		burns:     burns,
		limit:     ratelimit.New(rateLimit),
		dryRun:    dryRun,
	}

	return &p
}

func (p *ParsingConsumer) Consume(delivery rmq.Delivery) {

	log := p.log

	payload := []byte(delivery.Payload())
	err := delivery.Ack()
	if err != nil {
		log.Error().Err(err).Msg("could not acknowledge delivery")
		return
	}

	var parsing jobs.Parsing
	err = json.Unmarshal(payload, &parsing)
	if err != nil {
		log.Error().Err(err).Msg("could not decode payload")
		return
	}

	log = log.With().
		Str("chain_id", parsing.ChainID).
		Strs("addresses", parsing.Addresses).
		Strs("events_types", parsing.EventTypes).
		Uint64("start_height", parsing.StartHeight).
		Uint64("end_height", parsing.EndHeight).
		Logger()

	err = p.parsings.UpdateStatus(jobs.StatusProcessing, parsing.ID)
	if err != nil {
		p.log.Error().Err(err).Msg("could not update parsing job status")
		return
	}

	notify := func(err error, dur time.Duration) {
		log.Error().Err(err).Dur("duration", dur).Msg("could not complete lambda invocation")
	}

	err = p.process(notify, payload)
	if err != nil {
		log.Error().Err(err).Msg("could not process parsing job")
		err = p.parsings.UpdateStatus(jobs.StatusFailed, parsing.ID)
	} else {
		log.Info().Msg("parsing job successfully processed")
		err = p.parsings.UpdateStatus(jobs.StatusFinished, parsing.ID)
	}

	if err != nil {
		log.Error().Err(err).Msg("could not update parsing job status")
		return
	}
}

func (p *ParsingConsumer) process(notify func(error, time.Duration), input []byte) error {

	if p.dryRun {
		return nil
	}

	var output []byte
	err := backoff.RetryNotify(func() error {

		p.limit.Take()

		input := &lambda.InvokeInput{
			FunctionName: aws.String("parsing_worker"),
			Payload:      input,
		}
		result, err := p.lambda.Invoke(input)
		var reqErr *lambda.TooManyRequestsException

		// retry if we ran out of concurrent lambdas
		if errors.As(err, &reqErr) {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// retry if we ran out of requests on the Ethereum API
		if strings.Contains(err.Error(), "Too Many Requests") {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// don't retry on any other error, for now
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not execute lambda: %w", err))
		}

		output = result.Payload
		return nil
	}, backoff.NewExponentialBackOff(), notify)
	if err != nil {
		return fmt.Errorf("could not successfully invoke parser: %w", err)
	}

	var result results.Parsing
	err = json.Unmarshal(output, &result)
	if err != nil {
		return fmt.Errorf("could not decode parsing result: %w", err)
	}

	err = p.mints.Upsert(result.Mints...)
	if err != nil {
		return fmt.Errorf("could not upsert mints: %w", err)
	}

	err = p.transfers.Upsert(result.Transfers...)
	if err != nil {
		return fmt.Errorf("could not upsert transfers: %w", err)
	}

	err = p.sales.Upsert(result.Sales...)
	if err != nil {
		return fmt.Errorf("could not upsert sales: %w", err)
	}

	err = p.burns.Upsert(result.Burns...)
	if err != nil {
		return fmt.Errorf("could not upsert burns: %w", err)
	}

	for _, mint := range result.Mints {
		action := jobs.Action{
			ID:        uuid.New().String(),
			NetworkID: "", // TODO
			Address:   "", // TODO
			TokenID:   mint.TokenID,
		}
	}

	return nil
}
