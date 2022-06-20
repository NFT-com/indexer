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
	"github.com/NFT-com/indexer/models/results"
)

type ParsingStage struct {
	ctx       context.Context
	log       zerolog.Logger
	lambda    *lambda.Client
	name      string
	transfers TransferStore
	sales     SaleStore
	owners    OwnerStore
	pub       BatchPublisher
	limit     ratelimit.Limiter
	dryRun    bool
}

func NewParsingStage(
	ctx context.Context,
	log zerolog.Logger,
	lambda *lambda.Client,
	name string,
	transfers TransferStore,
	sales SaleStore,
	owners OwnerStore,
	pub BatchPublisher,
	limit ratelimit.Limiter,
	dryRun bool,
) *ParsingStage {

	p := ParsingStage{
		ctx:       ctx,
		log:       log,
		lambda:    lambda,
		name:      name,
		transfers: transfers,
		sales:     sales,
		owners:    owners,
		pub:       pub,
		limit:     limit,
		dryRun:    dryRun,
	}

	return &p
}

func (p *ParsingStage) HandleMessage(m *nsq.Message) error {

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

func (p *ParsingStage) process(payload []byte) error {

	if p.dryRun {
		return nil
	}

	p.limit.Take()

	// invoke the lambda and retry on retriable errors
	input := &lambda.InvokeInput{
		FunctionName: aws.String(p.name),
		Payload:      payload,
	}
	output, err := p.lambda.Invoke(p.ctx, input)
	if err != nil {
		return fmt.Errorf("could not invoke lambda: %w", err)
	}

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

	var result results.Parsing
	err = json.Unmarshal(output.Payload, &result)
	if err != nil {
		return fmt.Errorf("could not decode parsing result: %w", err)
	}

	payloads := make([][]byte, 0, len(result.Additions))
	for _, addition := range result.Additions {
		payload, err := json.Marshal(addition)
		if err != nil {
			return fmt.Errorf("could not encode addition job: %w", err)
		}
		payloads = append(payloads, payload)
	}

	err = p.pub.MultiPublish(params.TopicAddition, payloads)
	if err != nil {
		return fmt.Errorf("could not publish addition jobs: %w", err)
	}

	err = p.transfers.Upsert(result.Transfers...)
	if err != nil {
		return fmt.Errorf("could not upsert transfers: %w", err)
	}

	err = p.sales.Upsert(result.Sales...)
	if err != nil {
		return fmt.Errorf("could not upsert sales: %w", err)
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
