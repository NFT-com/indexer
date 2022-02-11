package ethereum_test

import (
	"context"
	"errors"
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
	subscription := mocks.BaselineSubscription(t)

	tests := []struct {
		name        string
		log         zerolog.Logger
		client      ethereum.Client
		network     string
		chain       string
		assertValue assert.ValueAssertionFunc
		assertError assert.ErrorAssertionFunc
	}{
		{
			name:        "should return error on missing client",
			log:         zerolog.Nop(),
			client:      nil,
			network:     "ethereum",
			chain:       "mainnet",
			assertValue: assert.Nil,
			assertError: assert.Error,
		},
		{
			name:        "should return error on missing network",
			log:         zerolog.Nop(),
			client:      mocks.BaselineClient(t, subscription),
			network:     "",
			chain:       "mainnet",
			assertValue: assert.Nil,
			assertError: assert.Error,
		},
		{
			name:        "should return error on missing network",
			log:         zerolog.Nop(),
			client:      mocks.BaselineClient(t, subscription),
			network:     "ethereum",
			chain:       "",
			assertValue: assert.Nil,
			assertError: assert.Error,
		},
		{
			name:        "should return parser correctly",
			log:         zerolog.Nop(),
			client:      mocks.BaselineClient(t, subscription),
			network:     "ethereum",
			chain:       "mainnet",
			assertValue: assert.NotNil,
			assertError: assert.Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			subs, err := ethereum.NewParser(test.log, test.client, test.network, test.chain)
			test.assertError(t, err)
			test.assertValue(t, subs)
		})
	}
}

func TestParser_Parse(t *testing.T) {
	var (
		ctx        = context.Background()
		log        = zerolog.Nop()
		mockClient = mocks.BaselineClient(t, nil)
		network    = "ethereum"
		chain      = "mainnet"
		b          = block.Block("block_1")
	)

	t.Run("return error filtering log", func(t *testing.T) {
		parser, err := ethereum.NewParser(log, mockClient, network, chain)
		require.NoError(t, err)

		mockClient.FilterLogsFunc = func(context.Context, goethereum.FilterQuery) ([]types.Log, error) {
			return nil, errors.New("failed to filter logs")
		}

		events, err := parser.Parse(ctx, &b)
		assert.Error(t, err)
		assert.Nil(t, events)
	})

	t.Run("should parse event correctly", func(t *testing.T) {
		parser, err := ethereum.NewParser(log, mockClient, network, chain)
		require.NoError(t, err)

		mockClient.FilterLogsFunc = func(_ context.Context, q goethereum.FilterQuery) ([]types.Log, error) {
			if q.BlockHash == nil && q.BlockHash.String() != "0x000000000000000000000000000000000000000000000000000000000000000b" {
				return nil, errors.New("bad block hash")
			}

			return mocks.GenericEthereumLogs, nil
		}

		events, err := parser.Parse(ctx, &b)
		assert.Error(t, err)
		assert.Len(t, events, 3)
	})
}
