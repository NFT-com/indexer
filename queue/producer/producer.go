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
}

func NewProducer(connection rmq.Connection, discoveryQueue string, parsingQueue string) (*Producer, error) {
	p := Producer{
		connection:     connection,
		discoveryQueue: discoveryQueue,
		parsingQueue:   parsingQueue,
	}

	return &p, nil
}

func (p *Producer) PublishDiscoveryJob(discoveryJob jobs.Discovery) error {
	q, err := p.connection.OpenQueue(p.discoveryQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	jobPayload, err := json.Marshal(discoveryJob)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	err = q.PublishBytes(jobPayload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}

func (p *Producer) PublishParsingJob(parsingJob jobs.Parsing) error {
	q, err := p.connection.OpenQueue(p.parsingQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	jobPayload, err := json.Marshal(parsingJob)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	err = q.PublishBytes(jobPayload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}
