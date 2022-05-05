package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/models/jobs"
)

type Producer struct {
	connection   rmq.Connection
	parsingQueue string
	actionQueue  string
}

func NewProducer(connection rmq.Connection, parsingQueue string, actionQueue string) (*Producer, error) {

	p := Producer{
		connection:   connection,
		parsingQueue: parsingQueue,
		actionQueue:  actionQueue,
	}

	return &p, nil
}

func (p *Producer) PublishParsingJob(job *jobs.Parsing) error {

	q, err := p.connection.OpenQueue(p.parsingQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	err = q.PublishBytes(payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}

func (p *Producer) PublishActionJob(job *jobs.Action) error {

	q, err := p.connection.OpenQueue(p.actionQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	err = q.PublishBytes(payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}
