package watcher

import (
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/job"
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

func (j *Job) Watch(discoveryJobs chan job.Discovery, parsingJobs chan job.Parsing) error {
	for {
		select {
		case newJob := <-discoveryJobs:
			err := j.publishDiscoveryJob(newJob)
			if err != nil {
				return err
			}
		case newJob := <-parsingJobs:
			err := j.publishParsingJob(newJob)
			if err != nil {
				return err
			}
		case <-j.close:
			return nil
		}
	}
}

func (j *Job) Close() {
	close(j.close)
}

func (j *Job) publishDiscoveryJob(newJob job.Discovery) error {
	err := j.messageProducer.PublishDiscoveryJob(newJob)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get publish discovery job")
		return err
	}

	err = j.apiClient.UpdateDiscoveryJobState(newJob.ID, job.StatusQueued)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get update discovery job status")
		return err
	}

	return nil
}

func (j *Job) publishParsingJob(newJob job.Parsing) error {
	err := j.messageProducer.PublishParsingJob(newJob)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get publish parsing job")
		return err
	}

	err = j.apiClient.UpdateParsingJobState(newJob.ID, job.StatusQueued)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get update parsing job status")
		return err
	}

	return err
}
