package mocks

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client struct {
	SubscribeNewHeadFunc func(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	HeaderByNumberFunc   func(ctx context.Context, number *big.Int) (*types.Header, error)
	FilterLogsFunc       func(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

func BaselineClient(t *testing.T, subscription ethereum.Subscription) *Client {
	t.Helper()

	c := Client{
		SubscribeNewHeadFunc: func(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
			return subscription, nil
		},
		HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*types.Header, error) {
			return GenericEthereumBlockHeader, nil
		},
		FilterLogsFunc: func(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
			return GenericEthereumLogs, nil
		},
	}

	return &c
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
