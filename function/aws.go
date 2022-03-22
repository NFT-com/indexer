package function

import (
	"errors"
	"fmt"

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

func (d *Client) Dispatch(functionName string, payload []byte) ([]byte, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      payload,
	}

	output, err := d.lambdaClient.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("could not invoke lambda: %w", err)
	}

	if output.StatusCode != nil && *output.StatusCode > 299 {
		if output.FunctionError != nil {
			return nil, fmt.Errorf("lambda finsihed with error: %s", *output.FunctionError)
		}

		return nil, fmt.Errorf("unexpected error during lambda invokation")
	}

	return output.Payload, nil
}
