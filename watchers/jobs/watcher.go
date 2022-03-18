package jobs

import (
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/client"
)

type Watcher struct {
	log             zerolog.Logger
	apiClient       *client.Client
	messageProducer *producer.Producer
	close           chan struct{}
}

func NewWatcher(log zerolog.Logger, apiClient *client.Client, messageProducer *producer.Producer) *Watcher {
	j := Watcher{
		log:             log.With().Str("component", "watchers").Logger(),
		apiClient:       apiClient,
		messageProducer: messageProducer,
		close:           make(chan struct{}),
	}

	return &j
}

func (j *Watcher) Watch(discoveryJobs chan jobs.Discovery, parsingJobs chan jobs.Parsing) error {
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

func (j *Watcher) Close() {
	close(j.close)
}

func (j *Watcher) publishDiscoveryJob(newJob jobs.Discovery) error {
	err := j.messageProducer.PublishDiscoveryJob(newJob)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get publish discovery job")
		return err
	}

	err = j.apiClient.UpdateDiscoveryJobState(newJob.ID, jobs.StatusQueued)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get update discovery job status")
		return err
	}

	return nil
}

func (j *Watcher) publishParsingJob(newJob jobs.Parsing) error {
	err := j.messageProducer.PublishParsingJob(newJob)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get publish parsing job")
		return err
	}

	err = j.apiClient.UpdateParsingJobState(newJob.ID, jobs.StatusQueued)
	if err != nil {
		j.log.Error().Err(err).Msg("could not get update parsing job status")
		return err
	}

	return err
}