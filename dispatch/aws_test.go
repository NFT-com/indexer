package dispatch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/testing/mocks"

	"github.com/aws/aws-sdk-go/service/lambda"
)

func TestNewClient(t *testing.T) {
	tts := []struct {
		name          string
		lambdaClient  dispatch.Lambda
		store         store.Storer
		expectedError bool
	}{
		{
			name:          "should return error on missing lambda client",
			lambdaClient:  nil,
			store:         mocks.BaselineStore(t),
			expectedError: true,
		},
		{
			name:          "should return error on missing store",
			lambdaClient:  mocks.BaselineLambda(t),
			store:         nil,
			expectedError: true,
		},
		{
			name:          "should return the correct client",
			lambdaClient:  mocks.BaselineLambda(t),
			store:         mocks.BaselineStore(t),
			expectedError: false,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			client, err := dispatch.NewClient(tt.lambdaClient, tt.store)
			if tt.expectedError && err == nil {
				t.Errorf("test %s failed expected error but got none", tt.name)
				return
			}

			if !tt.expectedError && client == nil {
				t.Errorf("test %s failed expected client but found none", tt.name)
				return
			}
		})
	}
}

func TestClient_Dispatch(t *testing.T) {
	t.Run("should return an error on failed store request", func(t *testing.T) {
		var (
			mockedLambdaClient = mocks.BaselineLambda(t)
			mockedStore        = mocks.BaselineStore(t)
		)

		client, err := dispatch.NewClient(mockedLambdaClient, mockedStore)
		if err != nil {
			t.Errorf("failed to create client")
			return
		}

		var (
			ctx = context.Background()
			e   = mocks.GenericEvents[0]
		)

		mockedStore.GetContractTypeFunc = func(_ context.Context, _, _, _ string) (string, error) {
			return "", errors.New("failed to get abi")
		}

		err = client.Dispatch(ctx, e)
		if err == nil {
			t.Errorf("expected an error but got none")
			return
		}
	})
	t.Run("should return an error on failed invoke", func(t *testing.T) {
		var (
			mockedLambdaClient = mocks.BaselineLambda(t)
			mockedStore        = mocks.BaselineStore(t)
		)

		client, err := dispatch.NewClient(mockedLambdaClient, mockedStore)
		if err != nil {
			t.Errorf("failed to create client")
			return
		}

		var (
			ctx = context.Background()
			e   = mocks.GenericEvents[0]
		)

		mockedStore.GetContractTypeFunc = func(_ context.Context, _, _, _ string) (string, error) {
			return "custom", nil
		}

		mockedLambdaClient.InvokeFunc = func(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
			return nil, errors.New("failed to invoke lambda")
		}

		err = client.Dispatch(ctx, e)
		if err == nil {
			t.Errorf("expected an error but got none")
			return
		}
	})
	t.Run("should return successfully invoke with function name", func(t *testing.T) {
		var (
			mockedLambdaClient = mocks.BaselineLambda(t)
			mockedStore        = mocks.BaselineStore(t)
		)

		client, err := dispatch.NewClient(mockedLambdaClient, mockedStore)
		if err != nil {
			t.Errorf("failed to create client")
			return
		}

		var (
			ctx = context.Background()
			e   = mocks.GenericEvents[0]
		)

		err = client.Dispatch(ctx, e)
		if err != nil {
			t.Errorf("expected an no error but got one")
			return
		}
	})
}
