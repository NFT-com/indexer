package consumer

import (
	"fmt"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/function"
)

type ParseConsumer struct {
	dispatch function.Dispatcher
}

func NewParseConsumer(dispatch function.Dispatcher) (*ParseConsumer, error) {
	c := ParseConsumer{
		dispatch: dispatch,
	}

	return &c, nil
}

func (d *ParseConsumer) Consume(delivery rmq.Delivery) {
	fmt.Println(delivery.Payload())
	
	if err := delivery.Ack(); err != nil {
		fmt.Println(5, err)
		// FIXME: Logg
	}
}
