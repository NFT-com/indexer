package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/queue"
)

type DiscoveryConsumer struct {
	dispatch function.Dispatcher
}

func NewDiscoveryConsumer(dispatch function.Dispatcher) (*DiscoveryConsumer, error) {
	c := DiscoveryConsumer{
		dispatch: dispatch,
	}

	return &c, nil
}

func (d *DiscoveryConsumer) Consume(delivery rmq.Delivery) {
	fmt.Println(delivery.Payload())

	payload := []byte(delivery.Payload())
	job := queue.DiscoveryJob{}

	err := json.Unmarshal(payload, &job)
	if err != nil {
		fmt.Println(1, err)
		// FIXME: Logg
		err = delivery.Reject()
		if err != nil {
			fmt.Println(2, err)
			// FIXME: LOG
		}
		return
	}

	err = d.dispatch.Dispatch(job.ChainType, payload)
	if err != nil {
		fmt.Println(3, err)
		// FIXME: Logg
		err = delivery.Reject()
		if err != nil {
			fmt.Println(4, err)
			// FIXME: LOG
		}
		return
	}

	err = delivery.Ack()
	if err != nil {
		fmt.Println(5, err)
		// FIXME: Logg
	}
}
