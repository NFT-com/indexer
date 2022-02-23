package mocks

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client struct {
	NetworkIDFunc        func(ctx context.Context) (*big.Int, error)
	ChainIDFunc          func(ctx context.Context) (*big.Int, error)
	SubscribeNewHeadFunc func(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	HeaderByNumberFunc   func(ctx context.Context, number *big.Int) (*types.Header, error)
	FilterLogsFunc       func(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

func BaselineClient(t *testing.T, subscription ethereum.Subscription) *Client {
	t.Helper()

	c := Client{
		NetworkIDFunc: func(ctx context.Context) (*big.Int, error) {
			return GenericNetworkID, nil
		},
		ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
			return GenericChainID, nil
		},
		SubscribeNewHeadFunc: func(context.Context, chan<- *types.Header) (ethereum.Subscription, error) {
			return subscription, nil
		},
		HeaderByNumberFunc: func(context.Context, *big.Int) (*types.Header, error) {
			return GenericEthereumBlockHeader, nil
		},
		FilterLogsFunc: func(context.Context, ethereum.FilterQuery) ([]types.Log, error) {
			return GenericEthereumLogs, nil
		},
	}

	return &c
}

func (c *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.NetworkIDFunc(ctx)
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	return c.ChainIDFunc(ctx)
}

func (c *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return c.SubscribeNewHeadFunc(ctx, ch)
}

func (c *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.HeaderByNumberFunc(ctx, number)
}

func (c *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return c.FilterLogsFunc(ctx, q)
}
