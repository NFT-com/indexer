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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
)

type ActionConsumer struct {
	log         zerolog.Logger
	lambda      *lambda.Lambda
	actions     ActionStore
	chains      ChainStore
	collections CollectionStore
	nfts        NFTStore
	traits      TraitStore
	dryRun      bool
}

func NewActionConsumer(
	log zerolog.Logger,
	lambda *lambda.Lambda,
	actions ActionStore,
	chains ChainStore,
	collections CollectionStore,
	nfts NFTStore,
	traits TraitStore,
	dryRun bool,
) *ActionConsumer {

	a := ActionConsumer{
		log:         log,
		lambda:      lambda,
		actions:     actions,
		chains:      chains,
		collections: collections,
		nfts:        nfts,
		traits:      traits,
		dryRun:      dryRun,
	}

	return &a
}

func (a *ActionConsumer) Consume(delivery rmq.Delivery) {

	log := a.log

	payload := []byte(delivery.Payload())
	err := delivery.Ack()
	if err != nil {
		log.Error().Err(err).Msg("could not acknowledge delivery")
		return
	}

	var action jobs.Action
	err = json.Unmarshal(payload, &action)
	if err != nil {
		log.Error().Err(err).Msg("could not decode payload")
		return
	}

	log = log.With().
		Str("chain_id", action.ChainID).
		Str("address", action.Address).
		Str("token_id", action.TokenID).
		Str("action_type", action.ActionType).
		Uint64("height", action.Height).
		Logger()

	err = a.actions.UpdateStatus(jobs.StatusProcessing, action.ID)
	if err != nil {
		a.log.Error().Err(err).Msg("could not update action job status")
		return
	}

	notify := func(err error, dur time.Duration) {
		log.Error().Err(err).Dur("duration", dur).Msg("could not complete lambda invocation")
	}

	switch action.ActionType {
	case jobs.ActionAddition:
		err = a.processAddition(notify, payload)
	case jobs.ActionOwnerChange:
		err = a.processOwnerChange(&action)
	}

	if err != nil {
		log.Error().Err(err).Msg("could not process action job")
		err = a.actions.UpdateStatus(jobs.StatusFailed, action.ID)
	} else {
		log.Info().Msg("action job successfully processed")
		err = a.actions.UpdateStatus(jobs.StatusFinished, action.ID)
	}

	if err != nil {
		log.Error().Err(err).Msg("could not update action job status")
		return
	}
}

func (a *ActionConsumer) processAddition(notify func(error, time.Duration), input []byte) error {

	if a.dryRun {
		return nil
	}

	var output []byte
	err := backoff.RetryNotify(func() error {

		input := &lambda.InvokeInput{
			FunctionName: aws.String("action_worker"),
			Payload:      input,
		}
		result, err := a.lambda.Invoke(input)
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

	var nft graph.NFT
	err = json.Unmarshal(output, &nft)
	if err != nil {
		return fmt.Errorf("could not decode NFT: %w", err)
	}

	collection, err := a.collections.RetrieveByAddress(chain.ID, nft.Contract, nft.ContractCollectionID)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	err = a.nfts.Upsert(&nft, collection.ID)
	if err != nil {
		return fmt.Errorf("could not store nft: %w", err)
	}

	for _, trait := range nft.Traits {
		err = a.traits.Upsert(trait)
		if err != nil {
			return fmt.Errorf("could not store trait: %w", err)
		}
	}

	return nil
}

func (a *ActionConsumer) processOwnerChange(action *jobs.Action) error {

	// TODO: check and sleep if there are any pending addition or owner change
	// jobs before this one for the same NFT

	var inputs inputs.OwnerChange
	err := json.Unmarshal(action.Data, &inputs)
	if err != nil {
		return fmt.Errorf("could not decode owner change inputs: %w", err)
	}

	err = a.nfts.ChangeOwner(inputs.NFTID, inputs.NewOwner)
	if err != nil {
		return fmt.Errorf("could not change owner: %w", err)
	}

	return nil
}
