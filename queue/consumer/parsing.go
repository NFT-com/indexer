package consumer

import (
	"encoding/json"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/jobs"
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
	parsingJob := job.Parsing{}

	err := json.Unmarshal(payload, &parsingJob)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to unmarshal message")
			return
		}

		d.log.Error().Err(err).Msg("failed to unmarshal message")
		return
	}

	retrievedParsingJob, err := d.apiClient.GetParsingJob(parsingJob.ID)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("failed to retrieve parsing job")
		return
	}

	if retrievedParsingJob.Status == jobs.StatusCanceled {
		err = delivery.Ack()
		if err != nil {
			d.log.Error().Err(err).Msg("failed to acknowledge message")
			return
		}
	}

	err = d.apiClient.UpdateParsingJobState(parsingJob.ID, jobs.StatusProcessing)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("failed to retrieve parsing job")
		return
	}

	status := jobs.StatusFinished
	err = d.dispatcher.Dispatch("parsing-85cd71d", payload)
	if err != nil {
		status = jobs.StatusFailed
		d.log.Error().Err(err).Msg("failed to dispatch message")
	}

	err = d.apiClient.UpdateParsingJobState(parsingJob.ID, job.Status(status))
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
