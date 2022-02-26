package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/NFT-com/indexer/queue"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"net/http"

	"github.com/adjust/rmq/v4"

	"github.com/NFT-com/indexer/function"
)

const (
	DefaultFunctionName = "default"
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

	payload := []byte(delivery.Payload())
	job := queue.ParseJob{}

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

	functions := []string{
		DefaultFunctionName,
		job.AddressType,
		job.Address,
	}

	err = d.dispatchFunctions(functions, payload)
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

func (d *ParseConsumer) dispatchFunctions(functions []string, payload []byte) error {
	for _, functionName := range functions {
		err := d.dispatch.Dispatch(functionName, payload)
		if err != nil {
			requestErr, ok := err.(awserr.RequestFailure)
			if ok && requestErr.StatusCode() == http.StatusNotFound {
				continue
			}

			return err
		}
	}

	return nil
}
