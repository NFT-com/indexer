package dispatch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name         string
		lambdaClient dispatch.Lambda
		store        store.Storer
		assertValue  assert.ValueAssertionFunc
		assertError  assert.ErrorAssertionFunc
	}{
		{
			name:         "return the correct client",
			lambdaClient: mocks.BaselineLambda(t),
			store:        mocks.BaselineStore(t),
			assertValue:  assert.NotNil,
			assertError:  assert.NoError,
		},
		{
			name:         "return error on missing lambda client",
			lambdaClient: nil,
			store:        mocks.BaselineStore(t),
			assertValue:  assert.Nil,
			assertError:  assert.Error,
		},
		{
			name:         "return error on missing store",
			lambdaClient: mocks.BaselineLambda(t),
			store:        nil,
			assertValue:  assert.Nil,
			assertError:  assert.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := dispatch.NewClient(test.lambdaClient, test.store)
			test.assertError(t, err)
			test.assertValue(t, client)
		})
	}
}

func TestClient_Dispatch(t *testing.T) {
	var (
		ctx                = context.Background()
		mockedLambdaClient = mocks.BaselineLambda(t)
		mockedStore        = mocks.BaselineStore(t)
		e                  = mocks.GenericEvents[0]
	)

	t.Run("return successfully invoke with function name", func(t *testing.T) {
		client, err := dispatch.NewClient(mockedLambdaClient, mockedStore)
		require.NoError(t, err)

		assert.Error(t, client.Dispatch(ctx, e))
	})

	t.Run("return an error on failed store request", func(t *testing.T) {
		client, err := dispatch.NewClient(mockedLambdaClient, mockedStore)
		require.NoError(t, err)

		mockedStore.GetContractTypeFunc = func(context.Context, string, string, string) (string, error) {
			return "", errors.New("failed to get abi")
		}

		assert.Error(t, client.Dispatch(ctx, e))
	})

	t.Run("return an error on failed invoke", func(t *testing.T) {
		client, err := dispatch.NewClient(mockedLambdaClient, mockedStore)
		require.NoError(t, err)

		mockedStore.GetContractTypeFunc = func(context.Context, string, string, string) (string, error) {
			return "custom", nil
		}

		mockedLambdaClient.InvokeFunc = func(*lambda.InvokeInput) (*lambda.InvokeOutput, error) {
			return nil, errors.New("failed to invoke lambda")
		}

		assert.Error(t, client.Dispatch(ctx, e))
	})
}
