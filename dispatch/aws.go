package dispatch

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/functions"
	"github.com/NFT-com/indexer/store"
)

const (
	customContractType = "custom"
)

type Lambda interface {
	Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

type Client struct {
	lambdaClient Lambda
	store        store.Storer
}

func NewClient(lambdaClient Lambda, store store.Storer) (*Client, error) {
	if lambdaClient == nil {
		return nil, errors.New("invalid lambda client")
	}
	if store == nil {
		return nil, errors.New("invalid store")
	}

	d := Client{
		lambdaClient: lambdaClient,
		store:        store,
	}

	return &d, nil
}

func (d *Client) Dispatch(ctx context.Context, e *event.Event) error {
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
