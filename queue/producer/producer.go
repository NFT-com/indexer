package producer

import (
	"encoding/json"

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
		return err
	}

	jobPayload, err := json.Marshal(discoveryJob)
	if err != nil {
		return err
	}

	err = q.PublishBytes(jobPayload)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) PublishParsingJob(parsingJob jobs.Parsing) error {
	q, err := p.connection.OpenQueue(p.parsingQueue)
	if err != nil {
		return err
	}

	jobPayload, err := json.Marshal(parsingJob)
	if err != nil {
		return err
	}

	err = q.PublishBytes(jobPayload)
	if err != nil {
		return err
	}

	return nil
}
