package web3

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Web3 struct {
	ethClient *ethclient.Client
	close     chan struct{}
}

func New(ctx context.Context, url string) (*Web3, error) {
	ethClient, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	w := Web3{
		ethClient: ethClient,
		close:     make(chan struct{}),
	}

	return &w, nil
}

func (w *Web3) ChainID(ctx context.Context) (string, error) {
	chainID, err := w.ethClient.ChainID(ctx)
	if err != nil {
		return "", err
	}

	return chainID.String(), nil
}

func (w *Web3) SubscribeToBlocks(ctx context.Context, blocks chan *big.Int) error {
	headerChannel := make(chan *types.Header)
	subscription, err := w.ethClient.SubscribeNewHead(ctx, headerChannel)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case header := <-headerChannel:
				blocks <- header.Number
			case <-w.close:
				subscription.Unsubscribe()
			}
		}
	}()

	return nil
}

func (w *Web3) GetLatestBlockHeight(ctx context.Context) (*big.Int, error) {
	header, err := w.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	return header.Number, nil
}

func (w *Web3) Close() {
	close(w.close)
	w.ethClient.Close()
}
