package pipeline

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/logs"
)

type ActionConsumer struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	actions       ActionStore
	chains        ChainStore
	collections   CollectionStore
	nfts          NFTStore
	traits        TraitStore
	limit         ratelimit.Limiter
	consumerQueue chan []byte
	close         chan struct{}
}

func NewActionConsumer(
	log zerolog.Logger,
	dispatcher function.Invoker,
	actions ActionStore,
	chains ChainStore,
	collections CollectionStore,
	nfts NFTStore,
	traits TraitStore,
	rateLimit int,
) *ActionConsumer {

	a := ActionConsumer{
		log:           log,
		dispatcher:    dispatcher,
		actions:       actions,
		chains:        chains,
		collections:   collections,
		nfts:          nfts,
		traits:        traits,
		limit:         ratelimit.New(rateLimit),
		consumerQueue: make(chan []byte, concurrentConsumers),
		close:         make(chan struct{}),
	}

	return &a
}

func (a *ActionConsumer) Consume(delivery rmq.Delivery) {

	a.log.Debug().Msg("received message to consume")

	payload := []byte(delivery.Payload())
	a.consumerQueue <- payload

	err := delivery.Ack()
	if err != nil {
		a.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (a *ActionConsumer) Run() {
	for i := 0; i < concurrentConsumers; i++ {
		go func() {
			for {
				select {
				case <-a.close:
					return
				case payload := <-a.consumerQueue:
					a.consume(payload)
				}
			}
		}()
	}
}

func (a *ActionConsumer) Close() {
	close(a.close)
}

func (a *ActionConsumer) consume(payload []byte) {

	var action jobs.Action
	err := json.Unmarshal(payload, &action)
	if err != nil {
		a.log.Error().Err(err).Msg("could not unmarshal message")
		return
	}

	err = a.actions.UpdateStatus(jobs.StatusProcessing, action.ID)
	if err != nil {
		a.handleError(action.ID, err, "could not update job status")
		return
	}

	a.log.Debug().
		Uint64("block", action.BlockNumber).
		Str("collection", action.Address).
		Str("standard", action.Standard).
		Str("event", action.Event).
		Str("token_id", action.TokenID).
		Str("action", action.Type).
		Msg("invoking function")

	name := actionName(&action)
	output, err := a.dispatcher.Invoke(name, payload)
	if err != nil {
		a.handleError(action.ID, err, "could not dispatch message")
		return
	}

	var nft graph.NFT
	err = json.Unmarshal(output, &nft)
	if err != nil {
		a.handleError(action.ID, err, "could not unmarshal output nft")
		return
	}

	err = a.processNFT(action.Type, &nft)
	if err != nil {
		a.handleError(action.ID, err, "could not process nft")
		return
	}

	err = a.actions.UpdateStatus(action.ID, jobs.StatusFinished)
	if err != nil {
		a.handleError(action.ID, err, "could not update job status")
		return
	}
}

func (a *ActionConsumer) handleError(actionID string, err error, message string) {
	updateErr := a.actions.UpdateStatus(actionID, jobs.StatusFailed)
	if updateErr != nil {
		a.log.Error().Err(updateErr).Msg("could not update job status")
	}

	a.log.Error().Err(err).Str("action_id", actionID).Msg(message)
}

func actionName(action *jobs.Action) string {
	h := sha256.New()

	s := strings.Join(
		[]string{
			"action",
			strings.ToLower(action.ChainType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}

func (a *ActionConsumer) processNFT(actionType string, nft *graph.NFT) error {

	chain, err := a.chains.Retrieve(nft.ChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}

	collection, err := a.collections.RetrieveByAddress(chain.ID, nft.Contract, nft.ContractCollectionID)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	switch actionType {

	case logs.ActionAddition.String():

		err = a.nfts.Upsert(nft, collection.ID)
		if err != nil {
			return fmt.Errorf("could not store nft: %w", err)
		}

		for _, trait := range nft.Traits {
			err = a.traits.Upsert(&trait)
			if err != nil {
				return fmt.Errorf("could not store trait: %w", err)
			}
		}

		a.log.Info().
			Str("collection", collection.ID).
			Str("nft", nft.ID).
			Str("name", nft.Name).
			Str("uri", nft.URI).
			Str("image", nft.Image).
			Int("traits", len(nft.Traits)).
			Msg("NFT details added")

	case logs.ActionOwnerChange.String():

		err = a.nfts.ChangeOwner(nft.ID, nft.Owner)
		if err != nil {
			return fmt.Errorf("could not update nft owner (nft %s): %w", nft.ID, err)
		}

		a.log.Info().
			Str("collection", collection.ID).
			Str("nft", nft.ID).
			Str("owner", nft.Owner).
			Msg("NFT owner updated")
	}

	return nil
}
