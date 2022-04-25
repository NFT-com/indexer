package action

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/service/client"
)

type Action struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	apiClient     *client.Client
	dataStore     Store
	jobCount      int
	consumerQueue chan []byte
	close         chan struct{}
}

func NewConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Invoker, dataStore Store, jobCount int) *Action {
	c := Action{
		log:           log,
		dispatcher:    dispatcher,
		apiClient:     apiClient,
		dataStore:     dataStore,
		jobCount:      jobCount,
		consumerQueue: make(chan []byte, jobCount),
		close:         make(chan struct{}),
	}

	return &c
}

func (d *Action) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())
	d.consumerQueue <- payload

	err := delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Action) Run() {
	for i := 0; i < d.jobCount; i++ {
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

	storedJob, err := d.apiClient.GetActionJob(job.ID)
	if err != nil {
		d.handleError(job.ID, err, "could not retrieve action job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		return
	}

	err = d.apiClient.UpdateActionJobStatus(job.ID, jobs.StatusProcessing)
	if err != nil {
		d.handleError(job.ID, err, "could not update job status")
		return
	}

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

	err = d.processNFT(job.ActionType, nft)
	if err != nil {
		d.handleError(job.ID, err, "could not process nft")
		return
	}

	err = d.apiClient.UpdateActionJobStatus(job.ID, jobs.StatusFinished)
	if err != nil {
		d.handleError(job.ID, err, "could not update job status")
		return
	}
}

func (d *Action) handleError(id string, err error, message string) {
	updateErr := d.apiClient.UpdateActionJobStatus(id, jobs.StatusFailed)
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
