package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

		input := &lambda.InvokeInput{
			FunctionName: aws.String(p.name),
			Payload:      payload,
		}
		result, err := p.lambda.Invoke(p.ctx, input)

		// retry if we ran out of Lambda requests
		if err != nil && strings.Contains(err.Error(), "Too Many Requests") {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// retry if we ran out of file handles
		if err != nil && strings.Contains(err.Error(), "too many opion files") {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// retry if we failed the DNS request
		if err != nil && strings.Contains(err.Error(), "no such host") {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// otherwise, fail permanently
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not invoke lambda: %w", err))
		}

		// if the function has not failed, we are done
		if result.FunctionError == nil {
			output = result.Payload
			return nil
		}

		// otherwise, process the function error
		var execErr results.Error
		err = json.Unmarshal(result.Payload, &execErr)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not decode error: %w", err))
		}

		// retry if we ran out of requests on the node
		if strings.Contains(execErr.Message, "Too Many Requests") {
			return fmt.Errorf("could not execute lambda: %s", execErr.Message)
		}

		// all other errors should be permanent for now
		return backoff.Permanent(fmt.Errorf("could not execute lambda: %s", execErr.Message))
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
