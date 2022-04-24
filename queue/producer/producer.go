package producer

import (
	"encoding/json"
	"fmt"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/jobs"
)

type Producer struct {
	connection     rmq.Connection
	discoveryQueue string
	parsingQueue   string
	additionQueue  string
}

func NewProducer(connection rmq.Connection, discoveryQueue string, parsingQueue string, additionQueue string) (*Producer, error) {
	p := Producer{
		connection:     connection,
		discoveryQueue: discoveryQueue,
		parsingQueue:   parsingQueue,
		additionQueue:  additionQueue,
	}

	return &p, nil
}

func (p *Producer) PublishDiscoveryJob(job *jobs.Discovery) error {
	q, err := p.connection.OpenQueue(p.discoveryQueue)
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

func (p *Producer) PublishAdditionJob(job *jobs.Addition) error {
	q, err := p.connection.OpenQueue(p.additionQueue)
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
