package ethereum_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewHistorical(t *testing.T) {
	var (
		ctx          = context.Background()
		log          = zerolog.Nop()
		client       = mocks.BaselineClient(t, nil)
		start  int64 = 1
		end    int64 = 6
	)

	t.Run("return correctly historical client", func(t *testing.T) {
		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		assert.Error(t, err)
		assert.Nil(t, historical)
	})

	t.Run("return error on failed header retrieval", func(t *testing.T) {
		end = 2

		historical, err := ethereum.NewHistorical(ctx, log, nil, start, end)
		assert.Error(t, err)
		assert.Nil(t, historical)
	})

	t.Run("return error on failed header retrieval", func(t *testing.T) {
		end = 2

		client.HeaderByNumberFunc = func(context.Context, *big.Int) (*types.Header, error) {
			return nil, errors.New("failed to retrieve header")
		}

		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		assert.Error(t, err)
		assert.Nil(t, historical)
	})
}

func TestHistoricalSource_Next(t *testing.T) {
	var (
		ctx          = context.Background()
		log          = zerolog.Nop()
		client       = mocks.BaselineClient(t, nil)
		start  int64 = 1
		end    int64 = 20
		newEnd int64 = 10
	)

	t.Run("return blocks correctly and stop on error", func(t *testing.T) {
		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
			h := types.Header{
				Number: big.NewInt(newEnd),
			}
			return &h, nil
		}

		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		require.Error(t, err)

		client.HeaderByNumberFunc = func(_ context.Context, number *big.Int) (*types.Header, error) {
			if number.Cmp(big.NewInt(10)) == 0 {
				return nil, errors.New("failed to get header")
			}

			h := types.Header{
				Root: common.HexToHash(number.String()),
			}
			return &h, nil
		}

		count := int64(0)
		for i := start; i <= end+1; i++ {
			b := historical.Next(ctx)

			if b != nil {
				count++
			}
		}

		assert.Equal(t, newEnd-start, count)
	})

	t.Run("return blocks correctly and stop on end", func(t *testing.T) {
		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
			h := types.Header{
				Number: big.NewInt(newEnd),
			}
			return &h, nil
		}

		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		require.NoError(t, err)

		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
			h := types.Header{
				Root: common.HexToHash(number.String()),
			}
			return &h, nil
		}

		count := int64(0)
		for i := start; i <= end+1; i++ {
			b := historical.Next(ctx)

			if b != nil {
				count++
			}
		}

		assert.Equal(t, (newEnd-start)+1, count)
	})
}

func TestHistoricalSource_Close(t *testing.T) {
	var (
		ctx          = context.Background()
		log          = zerolog.Nop()
		client       = mocks.BaselineClient(t, nil)
		start  int64 = 1
		end    int64 = 6
	)

	t.Run("return no error", func(t *testing.T) {
		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		require.Error(t, err)
		assert.NoError(t, historical.Close())
	})
}
