package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"

	"github.com/NFT-com/indexer/models/jobs"
)

type Producer struct {
	connection   *nsq.Producer
	parsingTopic string
	actionTopic  string
}

func NewProducer(connection *nsq.Producer, parsingTopic string, actionTopic string) (*Producer, error) {

	p := Producer{
		connection:   connection,
		parsingTopic: parsingTopic,
		actionTopic:  actionTopic,
	}

	return &p, nil
}

func (p *Producer) PublishParsingJob(job *jobs.Parsing) error {

	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	// Note: There is a possibility to MultiPublish (bulk publish)
	err = p.connection.Publish(p.parsingTopic, payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}

func (p *Producer) PublishActionJob(job *jobs.Action) error {

	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	// Note: There is a possibility to MultiPublish (bulk publish)
	err = p.connection.Publish(p.actionTopic, payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}
