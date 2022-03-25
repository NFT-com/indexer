package watcher

import (
	"fmt"

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

func NewJobWatcher(log zerolog.Logger, apiClient *client.Client, messageProducer *producer.Producer) *Job {
	j := Job{
		log:             log.With().Str("component", "watcher").Logger(),
		apiClient:       apiClient,
		messageProducer: messageProducer,
		close:           make(chan struct{}),
	}

	return &j
}

func (j *Job) Watch(discoveryJobs chan jobs.Discovery, parsingJobs chan jobs.Parsing) error {
	for {
		select {
		case newJob := <-discoveryJobs:
			err := j.publishDiscoveryJob(newJob)
			if err != nil {
				j.log.Error().
					Err(err).
					Str("id", newJob.ID).
					Str("block", newJob.BlockNumber).
					Str("status", string(newJob.Status)).
					Msg("could not publish discovery job")
				continue
			}
		case newJob := <-parsingJobs:
			err := j.publishParsingJob(newJob)
			if err != nil {
				j.log.Error().
					Err(err).
					Str("id", newJob.ID).
					Str("block", newJob.BlockNumber).
					Str("status", string(newJob.Status)).
					Msg("could not publish parsing job")
				continue
			}
		case <-j.close:
			return nil
		}
	}
}

func (j *Job) Close() {
	close(j.close)
}

func (j *Job) publishDiscoveryJob(newJob jobs.Discovery) error {
	err := j.messageProducer.PublishDiscoveryJob(newJob)
	if err != nil {
		return fmt.Errorf("could not get publish discovery job: %w", err)
	}

	err = j.apiClient.UpdateDiscoveryJobState(newJob.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update discovery job status: %w", err)
	}

	return nil
}

func (j *Job) publishParsingJob(newJob jobs.Parsing) error {
	err := j.messageProducer.PublishParsingJob(newJob)
	if err != nil {
		return fmt.Errorf("could not get publish parsing job: %w", err)
	}

	err = j.apiClient.UpdateParsingJobState(newJob.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	return nil
}
