package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"

	"github.com/NFT-com/indexer/models/jobs"
)

type JobCreator struct {
	producer     *nsq.Producer
	parsingTopic string
	actionTopic  string
}

func NewJobCreator(producer *nsq.Producer, parsingTopic string, actionTopic string) (*JobCreator, error) {

	j := JobCreator{
		producer:     producer,
		parsingTopic: parsingTopic,
		actionTopic:  actionTopic,
	}

	return &j, nil
}

func (j *JobCreator) PublishParsingJob(job *jobs.Parsing) error {

	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	// Note: There is a possibility to MultiPublish (bulk publish)
	err = j.producer.Publish(j.parsingTopic, payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}

func (j *JobCreator) PublishActionJob(job *jobs.Action) error {

	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	// Note: There is a possibility to MultiPublish (bulk publish)
	err = j.producer.Publish(j.actionTopic, payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}
