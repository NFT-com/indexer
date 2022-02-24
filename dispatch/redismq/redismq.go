package redismq

import (
	"encoding/json"

	"github.com/adjust/rmq/v4"
	
	"github.com/NFT-com/indexer/dispatch"
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

func (p *Producer) PublishDiscoveryJob(queueName string, job dispatch.DiscoveryJob) error {
	queue, err := p.connection.OpenQueue(queueName)
	if err != nil {
		return err
	}

	jobPayload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = queue.PublishBytes(jobPayload)
	if err != nil {
		return err
	}

	return nil
}
