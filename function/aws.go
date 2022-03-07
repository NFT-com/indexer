package function

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Lambda interface {
	Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

type Client struct {
	lambdaClient Lambda
}

func NewClient(lambdaClient Lambda) (*Client, error) {
	if lambdaClient == nil {
		return nil, errors.New("invalid lambda client")
	}

	d := Client{
		lambdaClient: lambdaClient,
	}

	return &d, nil
}

func (d *Client) Dispatch(functionName string, payload []byte) error {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      payload,
	}
	_, err := d.lambdaClient.Invoke(input)
	if err != nil {
		return err
	}

	return nil
}
