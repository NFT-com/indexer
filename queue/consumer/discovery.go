package consumer

import (
	"fmt"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/queue/producer"
)

type DiscoveryConsumer struct {
	producer *producer.Producer
}

func NewDiscoveryConsumer(producer *producer.Producer) (*DiscoveryConsumer, error) {
	c := DiscoveryConsumer{
		producer: producer,
	}

	return &c, nil
}

func (d *DiscoveryConsumer) Consume(delivery rmq.Delivery) {
	fmt.Println(delivery.Payload())
	_ = delivery.Ack()
}
