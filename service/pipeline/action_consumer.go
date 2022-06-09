package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
	"golang.org/x/crypto/sha3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/retry"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type ActionConsumer struct {
	ctx         context.Context
	log         zerolog.Logger
	client      *lambda.Client
	name        string
	actions     ActionStore
	collections CollectionStore
	nfts        NFTStore
	owners      OwnerStore
	traits      TraitStore
	limit       ratelimit.Limiter
	dryRun      bool
}

func NewActionConsumer(
	ctx context.Context,
	log zerolog.Logger,
	client *lambda.Client,
	name string,
	actions ActionStore,
	collections CollectionStore,
	nfts NFTStore,
	owners OwnerStore,
	traits TraitStore,
	limit ratelimit.Limiter,
	dryRun bool,
) *ActionConsumer {

	a := ActionConsumer{
		ctx:         ctx,
		log:         log,
		client:      client,
		name:        name,
		actions:     actions,
		collections: collections,
		nfts:        nfts,
		owners:      owners,
		traits:      traits,
		limit:       limit,
		dryRun:      dryRun,
	}

	return &a
}

func (a *ActionConsumer) Consume(delivery rmq.Delivery) {

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

	err = a.actions.Update(storage.One(action.ID), storage.SetStatus(jobs.StatusProcessing))
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
		err = a.actions.Update(storage.One(action.ID), storage.SetStatus(jobs.StatusFailed), storage.SetMessage(err.Error()))
	} else {
		log.Info().Msg("action job completed")
		err = a.actions.Update(storage.One(action.ID), storage.SetStatus(jobs.StatusFinished))
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
		log.Warn().Err(err).Dur("duration", dur).Msg("could not complete lambda invocation, retrying")
	}

	var output []byte
	err = backoff.RetryNotify(func() error {

		a.limit.Take()

		input := &lambda.InvokeInput{
			FunctionName: aws.String(a.name),
			Payload:      payload,
		}
		result, err := a.client.Invoke(a.ctx, input)

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
			return backoff.Permanent(fmt.Errorf("could not execute lambda: %s", execErr.Message))
		}

		output = result.Payload
		return nil
	}, retry.Indefinite(), notify)
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

	err = a.owners.AddCount(result.NFT.ID, result.NFT.Owner, int(result.NFT.Number))
	if err != nil {
		return fmt.Errorf("could not add owner count: %w", err)
	}

	return nil
}

func (a *ActionConsumer) processOwnerChange(action *jobs.Action) error {

	var inputs inputs.OwnerChange
	err := json.Unmarshal(action.InputData, &inputs)
	if err != nil {
		return fmt.Errorf("could not decode owner change inputs: %w", err)
	}

	collection, err := a.collections.One(action.ChainID, action.ContractAddress)
	if err != nil {
		return fmt.Errorf("could not retrieve collection: %w", err)
	}

	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s", action.ChainID, action.ContractAddress, action.TokenID)))
	nftID := uuid.Must(uuid.FromBytes(nftHash[:16]))

	err = a.nfts.Touch(nftID.String(), collection.ID, action.TokenID)
	if err != nil {
		return fmt.Errorf("could not touch NFT: %w", err)
	}

	err = a.owners.AddCount(nftID.String(), inputs.PrevOwner, -int(inputs.Number))
	if err != nil {
		return fmt.Errorf("could not decrease previous owner count: %w", err)
	}

	err = a.owners.AddCount(nftID.String(), inputs.NewOwner, int(inputs.Number))
	if err != nil {
		return fmt.Errorf("could not increase new owner count: %w", err)
	}

	return nil
}
