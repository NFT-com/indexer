package ethereum_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewHistorical(t *testing.T) {
	t.Run("should return error on failed header retrieval", func(t *testing.T) {
		var (
			ctx         = context.Background()
			log         = zerolog.Nop()
			start int64 = 1
			end   int64 = 2
		)

		historical, err := ethereum.NewHistorical(ctx, log, nil, start, end)
		if err == nil {
			t.Error("expected error creating historical but got none")
			return
		}

		if historical != nil {
			t.Error("unexpected historical struct")
			return
		}
	})
	t.Run("should return error on failed header retrieval", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			client       = mocks.BaselineClient(t, nil)
			start  int64 = 1
			end    int64 = 2
		)

		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
			return nil, errors.New("failed to retrieve header")
		}

		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		if err == nil {
			t.Error("expected error creating historical but got none")
			return
		}

		if historical != nil {
			t.Error("unexpected historical struct")
			return
		}
	})

	t.Run("should return correctly", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			client       = mocks.BaselineClient(t, nil)
			start  int64 = 1
			end    int64 = 6
		)
		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		if err != nil {
			t.Error("unexpected error creating historical")
			return
		}

		if historical == nil {
			t.Error("unexpected nil historical struct")
			return
		}
	})
}

func TestHistoricalSource_Next(t *testing.T) {
	t.Run("should return blocks correctly and stop on error", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			client       = mocks.BaselineClient(t, nil)
			start  int64 = 1
			end    int64 = 20
			newEnd int64 = 10
		)

		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
			h := types.Header{
				Number: big.NewInt(newEnd),
			}
			return &h, nil
		}

		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		if err != nil {
			t.Error("unexpected error creating historical")
			return
		}

		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
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

		if count != (newEnd - start) {
			t.Errorf("expected %v blocks but only got %v", newEnd-start, count)
			return
		}
	})
	t.Run("should return blocks correctly and stop on end", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			client       = mocks.BaselineClient(t, nil)
			start  int64 = 1
			end    int64 = 20
			newEnd int64 = 10
		)

		client.HeaderByNumberFunc = func(ctx context.Context, number *big.Int) (*types.Header, error) {
			h := types.Header{
				Number: big.NewInt(newEnd),
			}
			return &h, nil
		}

		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		if err != nil {
			t.Error("unexpected error creating historical")
			return
		}

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

		if count != (newEnd-start)+1 {
			t.Errorf("expected %v blocks but only got %v", newEnd-start+1, count)
			return
		}
	})
}

func TestHistoricalSource_Close(t *testing.T) {
	t.Run("should return no error", func(t *testing.T) {
		var (
			ctx          = context.Background()
			log          = zerolog.Nop()
			client       = mocks.BaselineClient(t, nil)
			start  int64 = 1
			end    int64 = 6
		)
		historical, err := ethereum.NewHistorical(ctx, log, client, start, end)
		if err != nil {
			t.Error("unexpected error creating historical")
			return
		}

		err = historical.Close()
		if err != nil {
			t.Error("unexpected error closing historical")
			return
		}
	})
}
