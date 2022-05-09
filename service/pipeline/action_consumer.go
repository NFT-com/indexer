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

	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type ActionConsumer struct {
	log         zerolog.Logger
	lambda      *lambda.Lambda
	name        string
	actions     ActionStore
	collections CollectionStore
	nfts        NFTStore
	traits      TraitStore
	limit       ratelimit.Limiter
	dryRun      bool
}

func NewActionConsumer(
	log zerolog.Logger,
	lambda *lambda.Lambda,
	name string,
	actions ActionStore,
	collections CollectionStore,
	nfts NFTStore,
	traits TraitStore,
	rateLimit uint,
	dryRun bool,
) *ActionConsumer {

	a := ActionConsumer{
		log:         log,
		lambda:      lambda,
		name:        name,
		actions:     actions,
		collections: collections,
		nfts:        nfts,
		traits:      traits,
		limit:       ratelimit.New(int(rateLimit)),
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

	err = a.process(payload)
	if err != nil {
		log.Error().Err(err).Msg("could not process payload")
		return
	}
}

func (a *ActionConsumer) process(payload []byte) error {

	var action jobs.Action
	err := json.Unmarshal(payload, &action)
	if err != nil {
		return fmt.Errorf("could not decode action job: %w", err)
	}
	log := a.log.With().
		Uint64("chain_id", action.ChainID).
		Str("contract_address", action.ContractAddress).
		Str("token_id", action.TokenID).
		Str("action_type", action.ActionType).
		Uint64("block_height", action.BlockHeight).
		Logger()

	err = a.actions.UpdateStatus(jobs.StatusProcessing, action.ID)
	if err != nil {
		return fmt.Errorf("could not update job status: %w", err)
	}

	switch action.ActionType {
	case jobs.ActionAddition:
		err = a.processAddition(payload, &action)
	case jobs.ActionOwnerChange:
		err = a.processOwnerChange(&action)
	default:
		err = fmt.Errorf("unknown action type (%s)", action.ActionType)
	}

	if err != nil {
		log.Error().Err(err).Msg("action job failed")
		err = a.actions.UpdateStatus(jobs.StatusFailed, action.ID)
	} else {
		log.Info().Msg("action job completed")
		err = a.actions.UpdateStatus(jobs.StatusFinished, action.ID)
	}

	if err != nil {
		return fmt.Errorf("could not update job status: %w", err)
	}

	return nil
}

func (a *ActionConsumer) processAddition(payload []byte, action *jobs.Action) error {

	if a.dryRun {
		return nil
	}

	collection, err := a.collections.One(action.ChainID, action.ContractAddress)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	notify := func(err error, dur time.Duration) {
		log.Error().Err(err).Dur("duration", dur).Msg("could not complete lambda invocation")
	}

	var output []byte
	err = backoff.RetryNotify(func() error {

		a.limit.Take()

		input := &lambda.InvokeInput{
			FunctionName: aws.String(a.name),
			Payload:      payload,
		}
		result, err := a.lambda.Invoke(input)
		var reqErr *lambda.TooManyRequestsException

		// retry if we ran out of concurrent lambdas
		if errors.As(err, &reqErr) {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// retry if we ran out of requests on the Ethereum API
		if err != nil && strings.Contains(err.Error(), "Too Many Requests") {
			return fmt.Errorf("could not invoke lambda: %w", err)
		}

		// don't retry on any other infrastructure error, for now
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not invoke lambda: %w", err))
		}

		// don't retry if the function failed for some reason
		if result.FunctionError != nil {
			var execErr results.Error
			err = json.Unmarshal(result.Payload, &execErr)
			if err != nil {
				return backoff.Permanent(fmt.Errorf("could not decode error: %w", err))
			}
			return backoff.Permanent(fmt.Errorf("could not execute lambda (%s): %s", execErr.Type, execErr.Message))
		}

		output = result.Payload
		return nil
	}, backoff.NewExponentialBackOff(), notify)
	if err != nil {
		return fmt.Errorf("could not successfully invoke parser: %w", err)
	}

	var result results.Addition
	err = json.Unmarshal(output, &result)
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

	return nil
}

func (a *ActionConsumer) processOwnerChange(action *jobs.Action) error {

	// TODO: check and sleep if there are any pending addition or owner change
	// jobs before this one for the same NFT

	var inputs inputs.OwnerChange
	err := json.Unmarshal(action.InputData, &inputs)
	if err != nil {
		return fmt.Errorf("could not decode owner change inputs: %w", err)
	}

	collection, err := a.collections.One(action.ChainID, action.ContractAddress)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	err = a.nfts.ChangeOwner(collection.ID, action.TokenID, inputs.NewOwner)
	if err != nil {
		return fmt.Errorf("could not change owner: %w", err)
	}

	return nil
}
