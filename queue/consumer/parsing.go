package consumer

import (
	"encoding/json"
	"log"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/client"
)

type Parsing struct {
	log        zerolog.Logger
	dispatcher function.Dispatcher
	apiClient  *client.Client
}

func NewParsingConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Dispatcher) (*Parsing, error) {
	c := Parsing{
		log:        log,
		dispatcher: dispatcher,
		apiClient:  apiClient,
	}

	return &c, nil
}

func (d *Parsing) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())
	var job jobs.Parsing

	err := json.Unmarshal(payload, &job)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to unmarshal message")
			return
		}

		d.log.Error().Err(err).Msg("failed to unmarshal message")
		return
	}

	storedJob, err := d.apiClient.GetParsingJob(job.ID)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("failed to retrieve parsing job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		err = delivery.Ack()
		if err != nil {
			d.log.Error().Err(err).Msg("failed to acknowledge message")
			return
		}
	}

	err = d.apiClient.UpdateParsingJobState(job.ID, jobs.StatusProcessing)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("failed to retrieve parsing job")
		return
	}

	status := jobs.StatusFinished
	lambdaOutput, err := d.dispatcher.Dispatch("parsing-85cd71d", payload)
	if err != nil {
		status = jobs.StatusFailed
		d.log.Error().Err(err).Msg("failed to dispatch message")
	}
	log.Println(lambdaOutput)

	err = d.apiClient.UpdateParsingJobState(job.ID, jobs.Status(status))
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("failed to retrieve parsing job")
		return
	}

	err = delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("failed to acknowledge message")
		return
	}
}
