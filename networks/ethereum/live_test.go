package ethereum_test

import (
	"context"
	"errors"
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
	t.Run("return the live source correctly", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
		)

		live, err := ethereum.NewLive(ctx, log, client)
		assert.NoError(t, err)
		assert.NotNil(t, live)
	})

	t.Run("return error on invalid client", func(t *testing.T) {
		var (
			ctx = context.Background()
			log = zerolog.Nop()
		)

		live, err := ethereum.NewLive(ctx, log, nil)
		assert.NoError(t, err)
		assert.NotNil(t, live)
	})

	t.Run("return error on failed to subscribe for headers", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
		)

		client.SubscribeNewHeadFunc = func(context.Context, chan<- *types.Header) (goethereum.Subscription, error) {
			return nil, errors.New("failed to subscribe to headers")
		}

		live, err := ethereum.NewLive(ctx, log, client)
		assert.NoError(t, err)
		assert.NotNil(t, live)
	})
}

func TestLiveSource_Next(t *testing.T) {
	t.Run("should return live block successfully", func(t *testing.T) {
		var (
			ctx           = context.Background()
			log           = zerolog.Nop()
			subscription  = mocks.BaselineSubscription(t)
			client        = mocks.BaselineClient(t, subscription)
			headerChannel = make(chan *types.Header)
		)

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
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
		)

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		go func() {
			<-time.After(time.Second)
			_ = live.Close()
		}()

		assert.Nil(t, live.Next(ctx))
	})

	t.Run("should cancel the context successfully", func(t *testing.T) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), time.Second)
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
		)

		defer cancel()

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		assert.Nil(t, live.Next(ctx))
	})

	t.Run("should close due to failed subscription with no error", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
			errChannel   = make(chan error)
		)

		subscription.ErrFunc = func() <-chan error {
			return errChannel
		}

		go func() {
			<-time.After(time.Second)
			errChannel <- nil
		}()

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		assert.Nil(t, live.Next(ctx))
	})

	t.Run("should close due to failed subscription with error", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
			errChannel   = make(chan error)
		)

		subscription.ErrFunc = func() <-chan error {
			return errChannel
		}

		go func() {
			<-time.After(time.Second)
			errChannel <- errors.New("failed to subscription")
		}()

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)

		assert.Nil(t, live.Next(ctx))
	})
}

func TestLiveSource_Close(t *testing.T) {
	t.Run("return no error", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			subscription = mocks.BaselineSubscription(t)
			client       = mocks.BaselineClient(t, subscription)
		)

		live, err := ethereum.NewLive(ctx, log, client)
		require.NoError(t, err)
		assert.Error(t, live.Close())
	})
}
