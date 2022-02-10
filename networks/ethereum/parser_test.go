package ethereum_test

import (
	"context"
	"errors"
	"testing"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewParser(t *testing.T) {
	subscription := mocks.BaselineSubscription(t)

	tts := []struct {
		name          string
		log           zerolog.Logger
		client        ethereum.Client
		network       string
		chain         string
		expectedError bool
	}{
		{
			name:          "should return error on missing client",
			log:           zerolog.Nop(),
			client:        nil,
			network:       "ethereum",
			chain:         "mainnet",
			expectedError: true,
		},
		{
			name:          "should return error on missing network",
			log:           zerolog.Nop(),
			client:        mocks.BaselineClient(t, subscription),
			network:       "",
			chain:         "mainnet",
			expectedError: true,
		},
		{
			name:          "should return error on missing network",
			log:           zerolog.Nop(),
			client:        mocks.BaselineClient(t, subscription),
			network:       "ethereum",
			chain:         "",
			expectedError: true,
		},
		{
			name:          "should return parser correctly",
			log:           zerolog.Nop(),
			client:        mocks.BaselineClient(t, subscription),
			network:       "ethereum",
			chain:         "mainnet",
			expectedError: false,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			subs, err := ethereum.NewParser(tt.log, tt.client, tt.network, tt.chain)
			if tt.expectedError && err == nil {
				t.Errorf("test %s failed expected error but got none", tt.name)
				return
			}

			if !tt.expectedError && subs == nil {
				t.Errorf("test %s failed expected subscriber but found none", tt.name)
				return
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	t.Run("should return error filtering log", func(t *testing.T) {
		var (
			log        = zerolog.Nop()
			mockClient = mocks.BaselineClient(t, nil)
			network    = "ethereum"
			chain      = "mainnet"
		)
		parser, err := ethereum.NewParser(log, mockClient, network, chain)
		if err != nil {
			t.Error("unexpected error creating parser")
			return
		}

		mockClient.FilterLogsFunc = func(ctx context.Context, q goethereum.FilterQuery) ([]types.Log, error) {
			return nil, errors.New("failed to filter logs")
		}

		ctx := context.Background()
		b := block.Block("block_1")
		_, err = parser.Parse(ctx, &b)
		if err == nil {
			t.Error("expected error parsing block")
			return
		}
	})

	t.Run("should parse event correctly", func(t *testing.T) {
		var (
			log        = zerolog.Nop()
			mockClient = mocks.BaselineClient(t, nil)
			network    = "ethereum"
			chain      = "mainnet"
		)
		parser, err := ethereum.NewParser(log, mockClient, network, chain)
		if err != nil {
			t.Error("unexpected error creating parser")
			return
		}

		mockClient.FilterLogsFunc = func(ctx context.Context, q goethereum.FilterQuery) ([]types.Log, error) {
			if q.BlockHash == nil && q.BlockHash.String() != "0x000000000000000000000000000000000000000000000000000000000000000b" {
				return nil, errors.New("bad block hash")
			}

			return mocks.GenericEthereumLogs, nil
		}

		ctx := context.Background()
		b := block.Block("block_1")
		events, err := parser.Parse(ctx, &b)
		if err != nil {
			t.Error("unexpected error parsing block")
			return
		}

		if len(events) != 3 {
			t.Error("unexpected length of events")
			return
		}
	})
}
