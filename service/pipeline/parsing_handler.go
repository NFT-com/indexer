package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/retry"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type ParsingHandler struct {
	ctx       context.Context
	log       zerolog.Logger
	lambda    *lambda.Client
	name      string
	parsings  ParsingStore
	actions   ActionStore
	transfers TransferStore
	sales     SaleStore
	limit     ratelimit.Limiter
	dryRun    bool
}

func NewParsingHandler(
	ctx context.Context,
	log zerolog.Logger,
	lambda *lambda.Client,
	name string,
	parsings ParsingStore,
	actions ActionStore,
	transfers TransferStore,
	sales SaleStore,
	limit ratelimit.Limiter,
	dryRun bool,
) *ParsingHandler {

	p := ParsingHandler{
		ctx:       ctx,
		log:       log,
		lambda:    lambda,
		name:      name,
		parsings:  parsings,
		actions:   actions,
		transfers: transfers,
		sales:     sales,
		limit:     limit,
		dryRun:    dryRun,
	}

	return &p
}

func (p *ParsingHandler) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	m.Finish()

	err := p.process(m.Body)
	if err != nil {
		log.Error().Err(err).Msg("could not process payload")
		return err
	}

	return nil
}

func (p *ParsingHandler) process(payload []byte) error {

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

	err = p.parsings.Update(storage.One(parsing.ID), storage.SetStatus(jobs.StatusProcessing))
	if err != nil {
		return fmt.Errorf("could not update job status: %w", err)
	}

	result, err := p.processParsing(payload)
	if err != nil {
		log.Error().Err(err).Msg("parsing job failed")
		err = p.parsings.Update(storage.One(parsing.ID), storage.SetStatus(jobs.StatusFailed), storage.SetMessage(err.Error()))
	} else {
		log.Info().
			Int("transfers", len(result.Transfers)).
			Int("sales", len(result.Sales)).
			Int("actions", len(result.Actions)).
			Msg("parsing job completed")
		err = p.parsings.Update(storage.One(parsing.ID), storage.SetStatus(jobs.StatusFinished))
	}

	if err != nil {
		return fmt.Errorf("could not update job status: %w", err)
	}

	return nil
}

func (p *ParsingHandler) processParsing(payload []byte) (*results.Parsing, error) {

	if p.dryRun {
		return &results.Parsing{}, nil
	}

	notify := func(err error, dur time.Duration) {
		log.Warn().Err(err).Dur("duration", dur).Msg("could not complete lambda invocation, retrying")
	}

	var output []byte
	err := backoff.RetryNotify(func() error {

		p.limit.Take()

		// invoke the lambda and retry on retriable errors
		input := &lambda.InvokeInput{
			FunctionName: aws.String(p.name),
			Payload:      payload,
		}
		result, err := p.lambda.Invoke(p.ctx, input)
		if results.Retriable(err) {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not invoke lambda: %w", err))
		}

		// check function error, if it is nil we are done
		if result.FunctionError == nil {
			output = result.Payload
			return nil
		}

		// otherwise, decode the function error
		var execErr results.Error
		err = json.Unmarshal(result.Payload, &execErr)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not decode error: %w", err))
		}
		if results.Retriable(execErr) {
			return fmt.Errorf("could not execute lambda: %w", execErr)
		}
		return backoff.Permanent(fmt.Errorf("could not execute lambda: %w", execErr))
	}, retry.Indefinite(), notify)
	if err != nil {
		return nil, fmt.Errorf("could not successfully invoke parser: %w", err)
	}

	var result results.Parsing
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("could not decode parsing result: %w", err)
	}

	err = p.transfers.Upsert(result.Transfers...)
	if err != nil {
		return nil, fmt.Errorf("could not upsert transfers: %w", err)
	}

	err = p.sales.Upsert(result.Sales...)
	if err != nil {
		return nil, fmt.Errorf("could not upsert sales: %w", err)
	}

	err = p.actions.Insert(result.Actions...)
	if err != nil {
		return nil, fmt.Errorf("could not upsert action jobs: %w", err)
	}

	// As we can't know in advance how many requests a Lambda will make, we will
	// wait here to take up any requests above one that we needed.
	for i := 1; i < int(result.Requests); i++ {
		p.limit.Take()
	}

	return &result, nil
}
