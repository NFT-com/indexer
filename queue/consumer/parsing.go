package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/NFT-com/indexer/job"
	"github.com/adjust/rmq/v4"
)

type Parsing struct {
}

func NewParsingConsumer() (*Parsing, error) {
	c := Parsing{}

	return &c, nil
}

func (d *Parsing) Consume(delivery rmq.Delivery) {
	fmt.Println(delivery.Payload())

	payload := []byte(delivery.Payload())
	parsingJob := job.Parsing{}

	err := json.Unmarshal(payload, &parsingJob)
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

	err = delivery.Ack()
	if err != nil {
		fmt.Println(5, err)
		// FIXME: Logg
	}
}
