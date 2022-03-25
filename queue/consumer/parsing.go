package consumer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/client"
)

type Parsing struct {
	log        zerolog.Logger
	dispatcher function.Dispatcher
	apiClient  *client.Client
	store      Store
}

func NewParsingConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Dispatcher, store Store) *Parsing {
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
		d.handleError(delivery, err, "could not unmarshal message")
		return
	}

	storedJob, err := d.apiClient.GetParsingJob(job.ID)
	if err != nil {
		d.handleError(delivery, err, "could not retrieve parsing job")
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
		d.handleError(delivery, err, "could not retrieve parsing job")
		return
	}

	name := functionName(job)
	lambdaOutput, err := d.dispatcher.Dispatch(name, payload)
	if err != nil {
		d.handleError(delivery, err, "could not dispatch message")
		err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFailed)
		if err != nil {
			d.handleError(delivery, err, "could not updating job state")
		}
		return
	}

	var jobResult []event.Event
	err = json.Unmarshal(lambdaOutput, &jobResult)
	if err != nil {
		d.handleError(delivery, err, "could not unmarshal job result")
		err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFailed)
		if err != nil {
			d.handleError(delivery, err, "could not updating job state")
		}
		return
	}

	err = d.processEvents(jobResult)
	if err != nil {
		d.handleError(delivery, err, "could not handle job result")
		err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFailed)
		if err != nil {
			d.handleError(delivery, err, "could not updating job state")
		}
		return
	}

	err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFinished)
	if err != nil {
		d.handleError(delivery, err, "could not updating job state")
		err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusFailed)
		if err != nil {
			d.handleError(delivery, err, "could not updating job state")
		}
		return
	}

	err = delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Parsing) processEvents(result []event.Event) error {
	for _, e := range result {
		err := d.store.InsertHistory(e)
		if err != nil {
			return fmt.Errorf("could not insert history: %w", err)
		}
	}

	return nil
}

func (d *Parsing) handleError(delivery rmq.Delivery, err error, message string) {
	log := d.log.Error()

	if err != nil {
		log = log.Err(err)
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

	s := strings.Join([]string{job.ChainType, job.StandardType, job.EventType}, "-")
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}
