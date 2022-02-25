package producer

import (
	"encoding/json"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/queue"
)

type Producer struct {
	connection rmq.Connection
}

func NewProducer(connection rmq.Connection) (*Producer, error) {
	p := Producer{
		connection: connection,
	}

	return &p, nil
}

func (p *Producer) PublishDiscoveryJob(queueName string, job queue.DiscoveryJob) error {
	q, err := p.connection.OpenQueue(queueName)
	if err != nil {
		return err
	}

	jobPayload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = q.PublishBytes(jobPayload)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) PublishParseJob(queueName string, job queue.ParseJob) error {
	q, err := p.connection.OpenQueue(queueName)
	if err != nil {
		return err
	}

	jobPayload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = q.PublishBytes(jobPayload)
	if err != nil {
		return err
	}

	return nil
}
