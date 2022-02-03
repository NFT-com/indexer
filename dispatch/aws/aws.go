package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/functions"
	"github.com/NFT-com/indexer/store"
)

const (
	customContractType = "custom"
)

type Dispatcher struct {
	lambdaClient *lambda.Lambda
	store        store.Storer
}

func New(lambdaClient *lambda.Lambda, store store.Storer) dispatch.Dispatcher {
	d := Dispatcher{
		lambdaClient: lambdaClient,
		store:        store,
	}

	return &d
}

func (d *Dispatcher) Dispatch(ctx context.Context, e *event.Event) error {
	contractType, err := d.store.GetContractType(ctx, e.Network, e.Chain, e.Address.Hex())
	if err != nil {
		return err
	}

	functionName := functions.Name(e.Network, e.Chain, contractType)
	if contractType == customContractType {
		functionName = functions.Name(e.Network, e.Chain, e.Address.Hex())
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      payload,
	}
	_, err = d.lambdaClient.Invoke(input)
	if err != nil {
		return err
	}

	return nil
}
