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
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
)

type ActionConsumer struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	jobStore      Store
	dataStore     Store
	limit         ratelimit.Limiter
	consumerQueue chan []byte
	close         chan struct{}
}

func NewActionConsumer(log zerolog.Logger, dispatcher function.Invoker, jobStore Store, dataStore Store, rateLimit int) *ActionConsumer {
	a := ActionConsumer{
		log:           log,
		dispatcher:    dispatcher,
		jobStore:      jobStore,
		dataStore:     dataStore,
		limit:         ratelimit.New(rateLimit),
		consumerQueue: make(chan []byte, concurrentConsumers),
		close:         make(chan struct{}),
	}

	return &a
}

func (a *ActionConsumer) Consume(delivery rmq.Delivery) {

	d.log.Debug().Msg("received message to consume")

	payload := []byte(delivery.Payload())
	d.consumerQueue <- payload

	err := delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (a *ActionConsumer) Run() {
	for i := 0; i < concurrentConsumers; i++ {
		go func() {
			for {
				select {
				case <-d.close:
					return
				case payload := <-d.consumerQueue:
					d.consume(payload)
				}
			}
		}()
	}
}

func (a *ActionConsumer) Close() {
	close(d.close)
}

func (a *ActionConsumer) consume(payload []byte) {

	var job jobs.ActionConsumer
	err := json.Unmarshal(payload, &job)
	if err != nil {
		d.log.Error().Err(err).Msg("could not unmarshal message")
		return
	}

	// job has been canceled meanwhile, no need to go further
	if job.Status != jobs.StatusCreated {
		return
	}

	storedJob, err := d.jobStore.ActionJob(job.ID)
	if err != nil {
		d.handleError(job.ID, err, "could not retrieve action job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		return
	}

	err = d.jobStore.UpdateActionJobStatus(job.ID, jobs.StatusProcessing)
	if err != nil {
		d.handleError(job.ID, err, "could not update job status")
		return
	}

	d.log.Debug().
		Uint64("block", job.BlockNumber).
		Str("collection", job.Address).
		Str("standard", job.Standard).
		Str("event", job.Event).
		Str("token_id", job.TokenID).
		Str("action", job.Type).
		Msg("invoking function")

	name := functionName(job)
	output, err := d.dispatcher.Invoke(name, payload)
	if err != nil {
		d.handleError(job.ID, err, "could not dispatch message")
		return
	}

	var nft chain.NFT
	err = json.Unmarshal(output, &nft)
	if err != nil {
		d.handleError(job.ID, err, "could not unmarshal output nft")
		return
	}

	err = d.processNFT(job.Type, nft)
	if err != nil {
		d.handleError(job.ID, err, "could not process nft")
		return
	}

	err = d.jobStore.UpdateActionJobStatus(job.ID, jobs.StatusFinished)
	if err != nil {
		d.handleError(job.ID, err, "could not update job status")
		return
	}
}

func (a *ActionConsumer) handleError(id string, err error, message string) {
	updateErr := d.jobStore.UpdateActionJobStatus(id, jobs.StatusFailed)
	if updateErr != nil {
		d.log.Error().Err(updateErr).Msg("could not update job status")
	}

	d.log.Error().Err(err).Str("job_id", id).Msg(message)
}

func functionName(job jobs.ActionConsumer) string {
	h := sha256.New()

	s := strings.Join(
		[]string{
			"action",
			strings.ToLower(job.ChainType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}

func (a *ActionConsumer) processNFT(actionType string, nft chain.NFT) error {

	chain, err := d.dataStore.Chain(nft.ChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}

	collection, err := d.dataStore.Collection(chain.ID, nft.Contract, nft.ContractCollectionID)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	switch actionType {

	case log.Addition.String():

		err = d.dataStore.UpsertNFT(nft, collection.ID)
		if err != nil {
			return fmt.Errorf("could not store nft: %w", err)
		}

		for _, trait := range nft.Traits {
			err = d.dataStore.UpsertTrait(trait)
			if err != nil {
				return fmt.Errorf("could not store trait: %w", err)
			}
		}

		d.log.Info().
			Str("collection", collection.ID).
			Str("nft", nft.ID).
			Str("name", nft.Name).
			Str("uri", nft.URI).
			Str("image", nft.Image).
			Int("traits", len(nft.Traits)).
			Msg("NFT details added")

	case log.OwnerChange.String():

		err = d.dataStore.UpdateNFTOwner(collection.ID, nft.ID, nft.Owner)
		if err != nil {
			return fmt.Errorf("could not update nft owner (nft %s): %w", nft.ID, err)
		}

		d.log.Info().
			Str("collection", collection.ID).
			Str("nft", nft.ID).
			Str("owner", nft.Owner).
			Msg("NFT owner updated")
	}

	return nil
}
