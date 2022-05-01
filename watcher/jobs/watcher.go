package jobs

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/postgres"
)

type Job struct {
	log             zerolog.Logger
	messageProducer *producer.Producer
	store           *postgres.Store
	delay           time.Duration
	close           chan struct{}
}

func New(log zerolog.Logger, messageProducer *producer.Producer, store *postgres.Store, delay time.Duration) *Job {
	j := Job{
		log:             log.With().Str("component", "watcher").Logger(),
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
	go j.watchActions()
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
			if len(jobs) == 0 {
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
			if len(jobs) == 0 {
				continue
			}

			j.handleParsingJobs(jobs)

		case <-j.close:
			return
		}
	}
}

func (j *Job) watchActions() {

	for {
		select {

		case <-time.After(j.delay):

			jobs, err := j.store.ActionJobs(jobs.StatusCreated)
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve action jobs")
				continue
			}
			if len(jobs) == 0 {
				continue
			}

			j.handleActionJobs(jobs)

		case <-j.close:
			return
		}
	}
}

func (j *Job) handleDiscoveryJobs(jobs []*jobs.Discovery) {

	for _, job := range jobs {

		err := j.publishDiscoveryJob(job)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("id", job.ID).
				Uint64("block", job.BlockNumber).
				Str("status", string(job.Status)).
				Msg("could not publish discovery job")
			continue
		}
	}

	j.log.Info().Int("jobs", len(jobs)).Msg("queued discovery jobs")
}

func (j *Job) publishDiscoveryJob(job *jobs.Discovery) error {

	if job.Status != jobs.StatusCreated {
		return nil
	}

	err := j.messageProducer.PublishDiscoveryJob(job)
	if err != nil {
		return fmt.Errorf("could not get publish discovery job: %w", err)
	}

	err = j.store.UpdateDiscoveryJobStatus(job.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update discovery job status: %w", err)
	}

	return nil
}

func (j *Job) handleParsingJobs(jobs []*jobs.Parsing) {

	for _, job := range jobs {
		err := j.publishParsingJob(job)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("id", job.ID).
				Uint64("block", job.BlockNumber).
				Str("status", string(job.Status)).
				Msg("could not publish parsing job")
			continue
		}
	}

	j.log.Info().Int("jobs", len(jobs)).Msg("queued parsing jobs")
}

func (j *Job) publishParsingJob(job *jobs.Parsing) error {

	if job.Status != jobs.StatusCreated {
		return nil
	}

	err := j.messageProducer.PublishParsingJob(job)
	if err != nil {
		return fmt.Errorf("could not get publish parsing job: %w", err)
	}

	err = j.store.UpdateParsingJobStatus(job.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	return nil
}

func (j *Job) handleActionJobs(jobs []*jobs.Action) {

	for _, job := range jobs {
		err := j.publishActionJob(job)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("id", job.ID).
				Uint64("block", job.BlockNumber).
				Str("status", string(job.Status)).
				Msg("could not publish action job")
			continue
		}
	}

	j.log.Info().Int("jobs", len(jobs)).Msg("queued action jobs")
}

func (j *Job) publishActionJob(job *jobs.Action) error {

	if job.Status != jobs.StatusCreated {
		return nil
	}

	err := j.messageProducer.PublishActionJob(job)
	if err != nil {
		return fmt.Errorf("could not publish action job: %w", err)
	}

	err = j.store.UpdateActionJobStatus(job.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update action job status: %w", err)
	}

	return nil
}
