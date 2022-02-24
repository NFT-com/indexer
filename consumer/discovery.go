package consumer

import (
	"fmt"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/dispatch/redismq"
)

type DiscoveryConsumer struct {
	producer *redismq.Producer
}

func NewDiscoveryConsumer(producer *redismq.Producer) (*DiscoveryConsumer, error) {
	c := DiscoveryConsumer{
		producer: producer,
	}

	return &c, nil
}

func (d *DiscoveryConsumer) Consume(delivery rmq.Delivery) {
	fmt.Println(delivery.Payload())
	_ = delivery.Ack()
}
