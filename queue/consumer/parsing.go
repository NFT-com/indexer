package consumer

import (
	"encoding/json"

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

func NewParsingConsumer(log zerolog.Logger, apiClient *client.Client) *Parsing {
	c := Parsing{
		log:       log,
		apiClient: apiClient,
	}

	return &c
}

func (d *Parsing) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())

	var job jobs.Parsing
	err := json.Unmarshal(payload, &job)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("could not unmarshal message")
			return
		}

		d.log.Error().Err(err).Msg("could not unmarshal message")
		return
	}

	storedJob, err := d.apiClient.GetParsingJob(job.ID)
	if err != nil {
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("could not retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("could not retrieve parsing job")
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
		if rejectErr := delivery.Reject(); rejectErr != nil {
			d.log.Error().Err(err).AnErr("reject_error", rejectErr).Msg("could not retrieve parsing job")
			return
		}

		d.log.Error().Err(err).Msg("could not retrieve parsing job")
		return
	}

	d.log.Info().Bytes("payload", payload).Msg("expected function call")

	err = delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}
