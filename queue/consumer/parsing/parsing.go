package parsing

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/service/client"
)

type Parsing struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	apiClient     *client.Client
	store         Store
	jobCount      int
	consumerQueue chan []byte
	close         chan struct{}
}

func NewConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Invoker, store Store, jobCount int) *Parsing {
	c := Parsing{
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

func (d *Parsing) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())
	d.consumerQueue <- payload

	err := delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Parsing) Run() {
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

func (d *Parsing) Close() {
	close(d.close)
}

func (d *Parsing) consume(payload []byte) {
	var job jobs.Parsing
	err := json.Unmarshal(payload, &job)
	if err != nil {
		d.log.Error().Err(err).Msg("could not unmarshal message")
		return
	}

	// job has been canceled meanwhile, no need to go further
	if job.Status != jobs.StatusCreated {
		return
	}

	storedJob, err := d.apiClient.GetParsingJob(job.ID)
	if err != nil {
		d.handleError(job.ID, err, "could not retrieve parsing job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		return
	}

	err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusProcessing)
	if err != nil {
		d.handleError(job.ID, err, "could not updating job status")
		return
	}

	fmt.Println(job)
	name := functionName(job)
	output, err := d.dispatcher.Invoke(name, payload)
	if err != nil {
		d.handleError(job.ID, err, "could not dispatch message")
		return
	}

	var logs []log.Log
	err = json.Unmarshal(output, &logs)
	if err != nil {
		d.handleError(job.ID, err, "could not unmarshal output logs")
		return
	}

	err = d.processLogs(job, logs)
	if err != nil {
		d.handleError(job.ID, err, "could not handle output logs")
		return
	}

	err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFinished)
	if err != nil {
		d.handleError(job.ID, err, "could not updating job status")
		return
	}
}

func (d *Parsing) handleError(id string, err error, message string) {
	updateErr := d.apiClient.UpdateParsingJobStatus(id, jobs.StatusFailed)
	if updateErr != nil {
		d.log.Error().Err(updateErr).Msg("could not update job status")
	}

	d.log.Error().Err(err).Str("job_id", id).Msg(message)
}

func functionName(job jobs.Parsing) string {
	h := sha256.New()

	s := strings.Join(
		[]string{
			strings.ToLower(job.ChainType),
			strings.ToLower(job.StandardType),
			strings.ToLower(job.EventType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}
