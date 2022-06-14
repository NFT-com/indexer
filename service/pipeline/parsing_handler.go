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

	err := p.process(m.Body)
	if results.Retriable(err) {
		log.Warn().Err(err).Msg("could not process message")
		return err
	}
	if err != nil {
		log.Error().Err(err).Msg("could not process message")
		// TODO: insert the failure into the DB
		err = nil
	}
	if err != nil {
		log.Fatal().Err(err).Msg("could not persist failure")
		return err
	}

	log.Trace().Msg("message processed")

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

	p.limit.Take()

	// invoke the lambda and retry on retriable errors
	input := &lambda.InvokeInput{
		FunctionName: aws.String(p.name),
		Payload:      payload,
	}
	output, err := p.lambda.Invoke(p.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("could not invoke lambda: %w", err)
	}

	var execErr *results.Error
	if output.FunctionError != nil {
		err = json.Unmarshal(output.Payload, &execErr)
	}
	if err != nil {
		return nil, fmt.Errorf("could not decode execution error: %w", err)
	}
	if execErr != nil {
		return nil, fmt.Errorf("could not execute lambda: %w", err)
	}

	var result results.Parsing
	err = json.Unmarshal(output.Payload, &result)
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
