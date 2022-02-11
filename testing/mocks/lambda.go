package mocks

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
)

type Lambda struct {
	InvokeFunc func(input *lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

func BaselineLambda(t *testing.T) *Lambda {
	t.Helper()

	c := Lambda{
		InvokeFunc: func(*lambda.InvokeInput) (*lambda.InvokeOutput, error) {
			return &GenericLambdaInvokeOutput, nil
		},
	}

	return &c
}

func (s *Lambda) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	return s.InvokeFunc(input)
}
