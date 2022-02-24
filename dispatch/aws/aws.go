package aws

import (
	"context"
	"encoding/json"
	"github.com/NFT-com/indexer/functions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/store"
)

const (
	customContractType = "custom"
)

type Dispatcher struct {
	lambdaClient *lambda.Lambda
	store        store.Storer
}

func New(lambdaClient *lambda.Lambda) dispatch.Dispatcher {
	d := Dispatcher{
		lambdaClient: lambdaClient,
	}

	return &d
}

func (d *Dispatcher) Dispatch(ctx context.Context, b *block.Block) error {
	payload, err := json.Marshal(b)
	if err != nil {
		return err
	}

	functionName := functions.Name(b.ChainID, b.NetworkID)
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
