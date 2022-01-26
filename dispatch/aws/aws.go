package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/store"
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
	/*contractType, err := d.store.GetContractType(ctx, e.Network, e.Chain, e.Address.Hex())
	if err != nil {
		return err
	}
	*/
	functionName := "test" /*
		switch contractType {
		case "erc721":
			functionName = functions.Name(e.Network, e.Chain, "erc721")
		case "erc1155":
			functionName = functions.Name(e.Network, e.Chain, "erc1155")
		case "custom":
			functionName = functions.Name(e.Network, e.Chain, e.Address.Hex())
		}*/

	payload, err := json.Marshal(e)
	if err != nil {
		// LOG
		return err
	}

	_, err = d.lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      payload,
	})
	if err != nil {
		return err
	}

	return nil
}
