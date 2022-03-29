package consumer

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
	log        zerolog.Logger
	dispatcher function.Invoker
	apiClient  *client.Client
	store      Store
}

func NewParsingConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Invoker, store Store) *Parsing {
	c := Parsing{
		log:        log,
		dispatcher: dispatcher,
		apiClient:  apiClient,
		store:      store,
	}

	return &c
}

func (d *Parsing) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())
	var job jobs.Parsing

	err := json.Unmarshal(payload, &job)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not unmarshal message")
		return
	}

	storedJob, err := d.apiClient.GetParsingJob(job.ID)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not retrieve parsing job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		err = delivery.Ack()
		if err != nil {
			d.log.Error().Err(err).Msg("could not acknowledge message")
			return
		}
	}

	err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusProcessing)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not retrieve parsing job")
		return
	}

	name := functionName(job)
	output, err := d.dispatcher.Invoke(name, payload)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not dispatch message")
		return
	}

	var logs []log.Log
	err = json.Unmarshal(output, &logs)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not unmarshal output logs")
		return
	}

	err = d.processLogs(logs)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not handle output logs")
		return
	}

	err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFinished)
	if err != nil {
		d.handleError(job.ID, delivery, err, "could not updating job status")
		return
	}

	err = delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Parsing) handleError(id string, delivery rmq.Delivery, err error, message string) {
	log := d.log.Error().Err(err).Str("job_id", id)

	updateErr := d.apiClient.UpdateParsingJobStatus(id, jobs.StatusFailed)
	if updateErr != nil {
		d.log.Error().Err(updateErr).Msg("could not update job status")
	}

	// rejects the message from the consumer
	rejectErr := delivery.Reject()
	if rejectErr != nil {
		log = log.AnErr("reject_error", rejectErr)
	}

	log.Msg(message)
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
