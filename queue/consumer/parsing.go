package consumer

import (
	"encoding/json"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/events"
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

func NewParsingConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Dispatcher, store Store) (*Parsing, error) {
	c := Parsing{
		log:        log,
		dispatcher: dispatcher,
		apiClient:  apiClient,
		store:      store,
	}

	return &c, nil
}

func (d *Parsing) Consume(delivery rmq.Delivery) {
	payload := []byte(delivery.Payload())
	var job jobs.Parsing

	err := json.Unmarshal(payload, &job)
	if err != nil {
		d.HandleError(delivery, err, "could not unmarshal message")
		return
	}

	storedJob, err := d.apiClient.GetParsingJob(job.ID)
	if err != nil {
		d.HandleError(delivery, err, "could not retrieve parsing job")
		return
	}

	if storedJob.Status == jobs.StatusCanceled {
		err = delivery.Ack()
		if err != nil {
			d.log.Error().Err(err).Msg("could not acknowledge message")
			return
		}
	}

	err = d.apiClient.UpdateParsingJobState(job.ID, jobs.StatusProcessing)
	if err != nil {
		d.HandleError(delivery, err, "could not retrieve parsing job")
		return
	}

	lambdaOutput, err := d.dispatcher.Dispatch("parsing-85cd71d", payload)
	if err != nil {
		d.HandleError(delivery, err, "could not dispatch message")
		err = d.apiClient.UpdateParsingJobState(job.ID, jobs.StatusFailed)
		if err != nil {
			d.HandleError(delivery, err, "could not updating job state")
		}
		return
	}

	var jobResult jobs.ParsingResult
	err = json.Unmarshal(lambdaOutput, &jobResult)
	if err != nil {
		d.HandleError(delivery, err, "failed unmarshal job result")
		err = d.apiClient.UpdateParsingJobState(job.ID, jobs.StatusFailed)
		if err != nil {
			d.HandleError(delivery, err, "could not updating job state")
		}
		return
	}

	err = d.HandlerJobResult(jobResult)
	if err != nil {
		d.HandleError(delivery, err, "could not handle job result")
		return
	}

	err = d.apiClient.UpdateParsingJobState(job.ID, jobs.StatusFinished)
	if err != nil {
		d.HandleError(delivery, err, "could not updating job state")
		return
	}

	err = delivery.Ack()
	if err != nil {
		d.log.Error().Err(err).Msg("could not acknowledge message")
		return
	}
}

func (d *Parsing) HandlerJobResult(result jobs.ParsingResult) error {
	for _, event := range result.RawEvents {
		err := d.store.InsertRawEvent(event)
		if err != nil {
			return err
		}
	}

	for _, event := range result.ParsedEvents {
		switch event.Type {
		case events.EventTypeMint:
			err := d.store.InsertNewNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, event.ToAddress)
			if err != nil {
				return err
			}
		case events.EventTypeUpdate, events.EventTypeBurn:
			err := d.store.UpdateNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, event.ToAddress)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Parsing) HandleError(delivery rmq.Delivery, err error, message string) {
	if rejectErr := delivery.Reject(); rejectErr != nil {
		log := d.log.Error()

		if err != nil {
			log = log.Err(err)
		}

		log.AnErr("reject_error", rejectErr).Msg(message)
		return
	}

	if err != nil {
		d.log.Error().Err(err).Msg(message)
	}
}
