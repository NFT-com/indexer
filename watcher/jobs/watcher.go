package jobs

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/client"
)

type Job struct {
	log             zerolog.Logger
	apiClient       *client.Client
	messageProducer *producer.Producer
	close           chan struct{}
}

func New(log zerolog.Logger, apiClient *client.Client, messageProducer *producer.Producer) *Job {
	j := Job{
		log:             log.With().Str("component", "watcher").Logger(),
		apiClient:       apiClient,
		messageProducer: messageProducer,
		close:           make(chan struct{}),
	}

	return &j
}

func (j *Job) Watch(discoveryJobs chan []jobs.Discovery, parsingJobs chan []jobs.Parsing) {
	for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
		go j.watch(discoveryJobs, parsingJobs)
	}
}

func (j *Job) Close() {
	close(j.close)
}

func (j *Job) watch(discoveryJobs chan []jobs.Discovery, parsingJobs chan []jobs.Parsing) {
	for {
		select {
		case jobs := <-discoveryJobs:
			j.handleDiscoveryJobs(jobs)
		case jobs := <-parsingJobs:
			j.handleParsingJobs(jobs)
		case <-j.close:
			return
		}
	}
}

func (j *Job) handleDiscoveryJobs(jobsList []jobs.Discovery) {
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

func (j *Job) publishDiscoveryJob(job jobs.Discovery) error {
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

func (j *Job) handleParsingJobs(jobsList []jobs.Parsing) {
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

func (j *Job) publishParsingJob(job jobs.Parsing) error {
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
