package addition

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

type Addition struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	apiClient     *client.Client
	store         Store
	jobCount      int
	consumerQueue chan []byte
	close         chan struct{}
}

func NewConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Invoker, store Store, jobCount int) *Addition {
	c := Addition{
		log:           log,
		dispatcher:    dispatcher,
		apiClient:     apiClient,
		store:         store,
		jobCount:      jobCount,
		consumerQueue: make(chan []byte, jobCount),
		close:         make(chan struct{}),
	}

	return &c
}

func (d *Addition) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())
	d.consumerQueue <- payload

	err := delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Addition) Run() {
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

func (d *Addition) Close() {
	close(d.close)
}

func (d *Addition) consume(payload []byte) {
	var job jobs.Addition
	err := json.Unmarshal(payload, &job)
	if err != nil {
		d.log.Error().Err(err).Msg("could not unmarshal message")
		return
	}

	// job has been canceled meanwhile, no need to go further
	if job.Status != jobs.StatusCreated {
		return
	}

	storedJob, err := d.apiClient.GetAdditionJob(job.ID)
	if err != nil {
		d.handleError(job.ID, err, "could not retrieve addition job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		return
	}

	err = d.apiClient.UpdateAdditionJobStatus(job.ID, jobs.StatusProcessing)
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

	err = d.processNFT(nft)
	if err != nil {
		d.handleError(job.ID, err, "could not process nft")
		return
	}

	err = d.apiClient.UpdateAdditionJobStatus(job.ID, jobs.StatusFinished)
	if err != nil {
		d.handleError(job.ID, err, "could not update job status")
		return
	}
}

func (d *Addition) handleError(id string, err error, message string) {
	updateErr := d.apiClient.UpdateAdditionJobStatus(id, jobs.StatusFailed)
	if updateErr != nil {
		d.log.Error().Err(updateErr).Msg("could not update job status")
	}

	d.log.Error().Err(err).Str("job_id", id).Msg(message)
}

func functionName(job jobs.Addition) string {
	h := sha256.New()

	s := strings.Join(
		[]string{
			strings.ToLower(job.ChainType),
			strings.ToLower(job.StandardType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}
