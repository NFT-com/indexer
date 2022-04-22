package parsing

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/service/client"
)

const concurrentConsumers = 4096

type Parsing struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	apiClient     *client.Client
	eventStore    Store
	dataStore     Store
	rateLimiter   <-chan time.Time
	consumerQueue chan []byte
	close         chan struct{}
}

func NewConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Invoker, eventStore Store, dataStore Store, rateLimit int) *Parsing {
	limiter := time.Tick(time.Second / time.Duration(rateLimit))
	c := Parsing{
		log:           log,
		dispatcher:    dispatcher,
		apiClient:     apiClient,
		eventStore:    eventStore,
		dataStore:     dataStore,
		rateLimiter:   limiter,
		consumerQueue: make(chan []byte, concurrentConsumers),
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
	for i := 0; i < concurrentConsumers; i++ {
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

	// Wait for rate limiter to have available spots.
	<-d.rateLimiter

	name := functionName(job)

	notify := func(err error, dur time.Duration) {
		d.log.Error().
			Err(err).
			Dur("retry_in", dur).
			Str("name", name).
			Int("payload_len", len(payload)).
			Msg("count not invoke lambda")
	}
	var output []byte
	_ = backoff.RetryNotify(func() error {
		output, err = d.dispatcher.Invoke(name, payload)
		if err != nil {
			return err
		}
		return nil
	}, backoff.NewExponentialBackOff(), notify)

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
			"parsing",
			strings.ToLower(job.ChainType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}
