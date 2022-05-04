package invoker

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Lambda interface {
	Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

type LambdaError struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
}

type AWS struct {
	client Lambda
}

func NewAWS(client Lambda) (*AWS, error) {
	if client == nil {
		return nil, errors.New("invalid lambda client")
	}

	d := AWS{
		client: client,
	}

	return &d, nil
}

func (d *AWS) Invoke(functionName string, payload []byte) ([]byte, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      payload,
	}

	output, err := d.client.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("could not invoke lambda: %w", err)
	}

	if output.StatusCode != nil && *output.StatusCode > 299 {
		if output.FunctionError != nil {
			return nil, fmt.Errorf("error during lambda runtime (status: %d, error: %s)", *output.StatusCode, *output.FunctionError)
		}

		return nil, fmt.Errorf("unexpected status from lambda runtime: %d", *output.StatusCode)
	}

	var lambdaError LambdaError
	err = json.Unmarshal(output.Payload, &lambdaError)
	if err != nil && !isArray(output.Payload) {
		return nil, fmt.Errorf("could not unmarshal output body error: %w", err)
	}

	if lambdaError.ErrorMessage != "" {
		return nil, fmt.Errorf("got an error from the lambda function: %s (error_type: %s)", lambdaError.ErrorMessage, lambdaError.ErrorType)
	}

	return output.Payload, nil
}

func isArray(in []byte) bool {
	return in[0] == '['
}
