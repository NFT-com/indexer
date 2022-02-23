package ethereum_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewParser(t *testing.T) {
	var (
		ctx    = context.Background()
		subs   = mocks.BaselineSubscription(t)
		client = mocks.BaselineClient(t, subs)
	)

	t.Run("should return no error", func(t *testing.T) {
		parser, err := ethereum.NewParser(ctx, zerolog.Nop(), client)
		require.NoError(t, err)

		assert.NotNil(t, parser)
	})

	t.Run("should return error on missing client", func(t *testing.T) {
		parser, err := ethereum.NewParser(ctx, zerolog.Nop(), nil)
		assert.Error(t, err)
		assert.Nil(t, parser)
	})

	t.Run("should return error on failed retrieval of network id", func(t *testing.T) {
		client.NetworkIDFunc = func(ctx context.Context) (*big.Int, error) {
			return nil, mocks.GenericError
		}

		parser, err := ethereum.NewParser(ctx, zerolog.Nop(), client)
		assert.Error(t, err)
		assert.Nil(t, parser)
	})

	t.Run("should return error on failed retrieval of chain id", func(t *testing.T) {
		client.NetworkIDFunc = func(ctx context.Context) (*big.Int, error) {
			return mocks.GenericNetworkID, nil
		}

		client.ChainIDFunc = func(ctx context.Context) (*big.Int, error) {
			return nil, mocks.GenericError
		}

		parser, err := ethereum.NewParser(ctx, zerolog.Nop(), client)
		assert.Error(t, err)
		assert.Nil(t, parser)
	})
}

func TestParser_Parse(t *testing.T) {
	var (
		ctx        = context.Background()
		log        = zerolog.Nop()
		mockClient = mocks.BaselineClient(t, nil)
		b          = block.Block("block_1")
	)

	t.Run("return error filtering log", func(t *testing.T) {
		parser, err := ethereum.NewParser(ctx, log, mockClient)
		require.NoError(t, err)

		mockClient.FilterLogsFunc = func(context.Context, goethereum.FilterQuery) ([]types.Log, error) {
			return nil, mocks.GenericError
		}

		events, err := parser.Parse(ctx, &b)
		assert.Error(t, err)
		assert.Nil(t, events)
	})

	t.Run("should parse event correctly", func(t *testing.T) {
		parser, err := ethereum.NewParser(ctx, log, mockClient)
		require.NoError(t, err)

		mockClient.FilterLogsFunc = func(_ context.Context, q goethereum.FilterQuery) ([]types.Log, error) {
			if q.BlockHash == nil && q.BlockHash.String() != "0x000000000000000000000000000000000000000000000000000000000000000b" {
				return nil, errors.New("bad block hash")
			}

			return mocks.GenericEthereumLogs, nil
		}

		events, err := parser.Parse(ctx, &b)
		assert.NoError(t, err)
		assert.Len(t, events, 3)
	})
}
