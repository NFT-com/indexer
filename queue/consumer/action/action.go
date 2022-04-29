package action

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

const concurrentConsumers = 1000

type Action struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	jobStore      Store
	dataStore     Store
	limit         ratelimit.Limiter
	consumerQueue chan []byte
	close         chan struct{}
}

func NewConsumer(log zerolog.Logger, dispatcher function.Invoker, jobStore Store, dataStore Store, rateLimit int) *Action {
	c := Action{
		log:           log,
		dispatcher:    dispatcher,
		jobStore:      jobStore,
		dataStore:     dataStore,
		limit:         ratelimit.New(rateLimit),
		consumerQueue: make(chan []byte, concurrentConsumers),
		close:         make(chan struct{}),
	}

	return &c
}

func (d *Action) Consume(delivery rmq.Delivery) {

	d.log.Debug().Msg("received message to consume")

	payload := []byte(delivery.Payload())
	d.consumerQueue <- payload

	err := delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Action) Run() {
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

func (d *Action) Close() {
	close(d.close)
}

func (d *Action) consume(payload []byte) {

	var job jobs.Action
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
		Str("block", job.BlockNumber).
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

func (d *Action) handleError(id string, err error, message string) {
	updateErr := d.jobStore.UpdateActionJobStatus(id, jobs.StatusFailed)
	if updateErr != nil {
		d.log.Error().Err(updateErr).Msg("could not update job status")
	}

	d.log.Error().Err(err).Str("job_id", id).Msg(message)
}

func functionName(job jobs.Action) string {
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
