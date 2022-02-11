package ethereum_test

import (
	"context"
	"testing"
	"time"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewLive(t *testing.T) {
	var (
		ctx          = context.Background()
		log          = zerolog.Nop()
		subscription = mocks.BaselineSubscription(t)
		client       = mocks.BaselineClient(t, subscription)
	)

	t.Run("return the live source correctly", func(t *testing.T) {
		client.SubscribeNewHeadFunc = func(context.Context, chan<- *types.Header) (goethereum.Subscription, error) {
			return subscription, nil
		}

		live, err := ethereum.NewLive(ctx, log, client)
		assert.NoError(t, err)
		assert.NotNil(t, live)
	})

	t.Run("return error on invalid client", func(t *testing.T) {
		client.SubscribeNewHeadFunc = func(context.Context, chan<- *types.Header) (goethereum.Subscription, error) {
			return subscription, nil
		}

		live, err := ethereum.NewLive(ctx, log, nil)
		assert.Error(t, err)
		assert.Nil(t, live)
	})

	t.Run("return error on failed to subscribe for headers", func(t *testing.T) {
		client.SubscribeNewHeadFunc = func(context.Context, chan<- *types.Header) (goethereum.Subscription, error) {
			return nil, mocks.GenericError
		}

		live, err := ethereum.NewLive(ctx, log, client)
		assert.Error(t, err)
		assert.Nil(t, live)
	})
}

func TestLiveSource_Next(t *testing.T) {
	var (
		ctx           = context.Background()
		log           = zerolog.Nop()
		subscription  = mocks.BaselineSubscription(t)
		client        = mocks.BaselineClient(t, subscription)
		headerChannel = make(chan *types.Header)
	)

	t.Run("should return live block successfully", func(t *testing.T) {
		client.SubscribeNewHeadFunc = func(ctx context.Context, ch chan<- *types.Header) (goethereum.Subscription, error) {
			go func() {
				for {
					h := <-headerChannel
					ch <- h
				}
			}()

			return subscription, nil
		}

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		headerChannel <- mocks.GenericEthereumBlockHeader

		b := live.Next(ctx)
		assert.Equal(t, mocks.GenericEthereumBlockHeader.Hash().Hex(), b.String())
	})

	t.Run("should close successfully", func(t *testing.T) {
		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		go func() {
			_ = live.Close()
		}()

		assert.Nil(t, live.Next(ctx))
	})

	t.Run("should cancel the context successfully", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		assert.Nil(t, live.Next(ctx))
	})

	t.Run("should close due to failed subscription with no error", func(t *testing.T) {
		errChannel := make(chan error)
		subscription.ErrFunc = func() <-chan error {
			return errChannel
		}

		go func() {
			errChannel <- nil
		}()

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		assert.Nil(t, live.Next(ctx))
	})

	t.Run("should close due to failed subscription with error", func(t *testing.T) {
		errChannel := make(chan error)
		subscription.ErrFunc = func() <-chan error {
			return errChannel
		}

		go func() {
			errChannel <- mocks.GenericError
		}()

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		assert.Nil(t, live.Next(ctx))
	})
}

func TestLiveSource_Close(t *testing.T) {
	var (
		ctx          = context.Background()
		log          = zerolog.Nop()
		subscription = mocks.BaselineSubscription(t)
		client       = mocks.BaselineClient(t, subscription)
	)

	t.Run("return no error", func(t *testing.T) {
		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)
		assert.NoError(t, live.Close())
	})
}