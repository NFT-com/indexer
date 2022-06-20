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
		limit:       limit,
		dryRun:      dryRun,
	}

	return &a
}

func (a *AdditionStage) HandleMessage(m *nsq.Message) error {

	err := a.process(m.Body)
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

func (a *AdditionStage) process(payload []byte) error {

	if a.dryRun {
		return nil
	}

	var addition jobs.Addition
	err := json.Unmarshal(payload, &addition)
	if err != nil {
		return fmt.Errorf("could not decode addition job: %w", err)
	}

	collection, err := a.collections.One(addition.ChainID, addition.ContractAddress)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	a.limit.Take()

	// invoke the lambda and retry on retriable errors
	input := &lambda.InvokeInput{
		FunctionName: aws.String(a.name),
		Payload:      payload,
	}
	output, err := a.lambda.Invoke(a.ctx, input)
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

	var result results.Addition
	err = json.Unmarshal(output.Payload, &result)
	if err != nil {
		return fmt.Errorf("could not decode NFT: %w", err)
	}

	result.NFT.CollectionID = collection.ID
	err = a.nfts.Insert(result.NFT)
	if err != nil {
		return fmt.Errorf("could not insert NFT: %w", err)
	}

	err = a.traits.Insert(result.Traits...)
	if err != nil {
		return fmt.Errorf("could not insert traits: %w", err)
	}

	err = a.owners.Add(&addition)
	if err != nil {
		return fmt.Errorf("could not add owner: %w", err)
	}

	// As we can't know in advance how many requests a Lambda will make, we will
	// wait here to take as many slots on the rate limiter as were needed.
	for i := 0; i < int(result.Requests); i++ {
		a.limit.Take()
	}

	return nil
}
