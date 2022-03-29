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

func (p *Producer) PublishDiscoveryJob(job jobs.Discovery) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	q, err := p.connection.OpenQueue(p.discoveryQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	err = q.PublishBytes(payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}

func (p *Producer) PublishDiscoveryJobs(jobs []jobs.Discovery) error {
	payloads := make([][]byte, 0, len(jobs))

	for _, job := range jobs {
		payload, err := json.Marshal(job)
		if err != nil {
			return fmt.Errorf("could not marshal payload: %w", err)
		}

		payloads = append(payloads, payload)
	}

	q, err := p.connection.OpenQueue(p.discoveryQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	err = q.PublishBytes(payloads...)
	if err != nil {
		return fmt.Errorf("could not publish jobs: %w", err)
	}

	return nil
}

func (p *Producer) PublishParsingJobs(jobs []jobs.Parsing) error {
	payloads := make([][]byte, 0, len(jobs))

	for _, job := range jobs {
		payload, err := json.Marshal(job)
		if err != nil {
			return fmt.Errorf("could not marshal payload: %w", err)
		}

		payloads = append(payloads, payload)
	}

	q, err := p.connection.OpenQueue(p.parsingQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	err = q.PublishBytes(payloads...)
	if err != nil {
		return fmt.Errorf("could not publish jobs: %w", err)
	}

	return nil
}

func (p *Producer) PublishParsingJob(job jobs.Parsing) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	q, err := p.connection.OpenQueue(p.parsingQueue)
	if err != nil {
		return fmt.Errorf("could not open connection with queue: %w", err)
	}

	err = q.PublishBytes(payload)
	if err != nil {
		return fmt.Errorf("could not publish job: %w", err)
	}

	return nil
}
