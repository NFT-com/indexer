package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type ParsingConsumer struct {
	log       zerolog.Logger
	lambda    *lambda.Lambda
	name      string
	parsings  ParsingStore
	actions   ActionStore
	transfers TransferStore
	sales     SaleStore
	limit     ratelimit.Limiter
	dryRun    bool
}

func NewParsingConsumer(
	log zerolog.Logger,
	lambda *lambda.Lambda,
	name string,
	parsings ParsingStore,
	actions ActionStore,
	transfers TransferStore,
	sales SaleStore,
	rateLimit uint,
	dryRun bool,
) *ParsingConsumer {

	p := ParsingConsumer{
		log:       log,
		lambda:    lambda,
		name:      name,
		parsings:  parsings,
		actions:   actions,
		transfers: transfers,
		sales:     sales,
		limit:     ratelimit.New(int(rateLimit)),
		dryRun:    dryRun,
	}

	return &p
}

func (p *ParsingConsumer) Consume(delivery rmq.Delivery) {

	payload := []byte(delivery.Payload())
	err := delivery.Ack()
	if err != nil {
		log.Error().Err(err).Msg("could not acknowledge delivery")
		return
	}

	err = p.process(payload)
	if err != nil {
		log.Error().Err(err).Msg("could not process payload")
		return
	}
}

func (p *ParsingConsumer) process(payload []byte) error {

	var parsing jobs.Parsing
	err := json.Unmarshal(payload, &parsing)
	if err != nil {
		return fmt.Errorf("could not unmarshal parsing job: %w", err)
	}

	log := p.log.With().
		Uint64("chain_id", parsing.ChainID).
		Strs("contract_addresses", parsing.ContractAddresses).
		Strs("event_hashes", parsing.EventHashes).
		Uint64("start_height", parsing.StartHeight).
		Uint64("end_height", parsing.EndHeight).
		Logger()

	err = p.parsings.UpdateStatus(jobs.StatusProcessing, parsing.ID)
	if err != nil {
		return fmt.Errorf("could not update job status: %w", err)
	}

	err = p.processParsing(payload, &parsing)
	if err != nil {
		log.Error().Err(err).Msg("could not process parsing job")
		err = p.parsings.UpdateStatus(jobs.StatusFailed, parsing.ID)
	} else {
		log.Info().Msg("parsing job successfully processed")
		err = p.parsings.UpdateStatus(jobs.StatusFinished, parsing.ID)
	}

	if err != nil {
		return fmt.Errorf("could not update job status: %w", err)
	}

	return nil
}

func (p *ParsingConsumer) processParsing(payload []byte, parsing *jobs.Parsing) error {

	if p.dryRun {
		return nil
	}

	notify := func(err error, dur time.Duration) {
		log.Error().Err(err).Dur("duration", dur).Msg("could not complete lambda invocation")
	}

	var output []byte
	err := backoff.RetryNotify(func() error {

		p.limit.Take()

		input := &lambda.InvokeInput{
			FunctionName: aws.String(p.name),
			Payload:      payload,
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

	err = p.transfers.Upsert(result.Transfers...)
	if err != nil {
		return fmt.Errorf("could not upsert transfers: %w", err)
	}

	err = p.sales.Upsert(result.Sales...)
	if err != nil {
		return fmt.Errorf("could not upsert sales: %w", err)
	}

	err = p.actions.Insert(result.Actions...)
	if err != nil {
		return fmt.Errorf("could not upsert action jobs: %w", err)
	}

	return nil
}
