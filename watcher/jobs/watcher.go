package jobs

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/service/postgres"
)

type Job struct {
	log             zerolog.Logger
	apiClient       *client.Client
	messageProducer *producer.Producer
	store           *postgres.Store
	delay           time.Duration
	close           chan struct{}
}

func New(log zerolog.Logger, apiClient *client.Client, messageProducer *producer.Producer, store *postgres.Store, delay time.Duration) *Job {
	j := Job{
		log:             log.With().Str("component", "watcher").Logger(),
		apiClient:       apiClient,
		messageProducer: messageProducer,
		store:           store,
		delay:           delay,
		close:           make(chan struct{}),
	}

	return &j
}

func (j *Job) Watch() {
	go j.watchDiscoveries()
	go j.watchParsings()
	go j.watchAdditions()
}

func (j *Job) Close() {
	close(j.close)
}

func (j *Job) watchDiscoveries() {
	for {
		select {
		case <-time.After(j.delay):
			jobs, err := j.store.DiscoveryJobs(jobs.StatusCreated)
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve discovery jobs")
				continue
			}
			j.handleDiscoveryJobs(jobs)
		case <-j.close:
			return
		}
	}
}

func (j *Job) watchParsings() {
	for {
		select {
		case <-time.After(j.delay):
			jobs, err := j.store.ParsingJobs(jobs.StatusCreated)
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve parsing jobs")
				continue
			}
			j.handleParsingJobs(jobs)
		case <-j.close:
			return
		}
	}
}

func (j *Job) watchAdditions() {
	for {
		select {
		case <-time.After(j.delay):
			jobs, err := j.store.AdditionJobs(jobs.StatusCreated)
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve addition jobs")
				continue
			}
			j.handleAdditionJobs(jobs)
		case <-j.close:
			return
		}
	}
}

func (j *Job) handleDiscoveryJobs(jobsList []*jobs.Discovery) {
	for _, job := range jobsList {
		err := j.publishDiscoveryJob(job)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("id", job.ID).
				Str("block", job.BlockNumber).
				Str("status", string(job.Status)).
				Msg("could not publish discovery job")
			continue
		}
	}
}

func (j *Job) publishDiscoveryJob(job *jobs.Discovery) error {
	if job.Status != jobs.StatusCreated {
		return nil
	}

	err := j.messageProducer.PublishDiscoveryJob(job)
	if err != nil {
		return fmt.Errorf("could not get publish discovery job: %w", err)
	}

	err = j.apiClient.UpdateDiscoveryJobStatus(job.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update discovery job status: %w", err)
	}

	return nil
}

func (j *Job) handleParsingJobs(jobsList []*jobs.Parsing) {
	for _, job := range jobsList {
		err := j.publishParsingJob(job)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("id", job.ID).
				Str("block", job.BlockNumber).
				Str("status", string(job.Status)).
				Msg("could not publish parsing job")
			continue
		}
	}
}

func (j *Job) publishParsingJob(job *jobs.Parsing) error {
	if job.Status != jobs.StatusCreated {
		return nil
	}

	err := j.messageProducer.PublishParsingJob(job)
	if err != nil {
		return fmt.Errorf("could not get publish parsing job: %w", err)
	}

	err = j.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	return nil
}

func (j *Job) handleAdditionJobs(jobsList []*jobs.Addition) {
	for _, job := range jobsList {
		err := j.publishAdditionJob(job)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("id", job.ID).
				Str("block", job.BlockNumber).
				Str("status", string(job.Status)).
				Msg("could not publish addition job")
			continue
		}
	}
}

func (j *Job) publishAdditionJob(job *jobs.Addition) error {
	if job.Status != jobs.StatusCreated {
		return nil
	}

	err := j.messageProducer.PublishAdditionJob(job)
	if err != nil {
		return fmt.Errorf("could not publish addition job: %w", err)
	}

	err = j.apiClient.UpdateAdditionJobStatus(job.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update addition job status: %w", err)
	}

	return nil
}
