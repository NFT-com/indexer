package consumer

import (
	"encoding/json"
	"github.com/NFT-com/indexer/service/client"
	"log"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/job"
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
	log.Println(string(delivery.Payload()))

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

	if retrievedParsingJob.Status == job.StatusCanceled {
		err = delivery.Ack()
		if err != nil {
			d.log.Error().Err(err).Msg("failed to acknowledge message")
			return
		}
	}

	err = d.dispatcher.Dispatch("test", payload)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("failed to dispatch message")
			return
		}

		d.log.Error().Err(err).Msg("failed to dispatch message")
		return
	}

	err = delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("failed to acknowledge message")
		return
	}
}