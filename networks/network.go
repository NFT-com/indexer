package networks

import (
	"context"
	"math/big"

	"github.com/NFT-com/indexer/log"
)

type Network interface {
	ChainID(ctx context.Context) (string, error)
	SubscribeToBlocks(ctx context.Context, blocks chan *big.Int) error
	GetLatestBlockHeight(ctx context.Context) (*big.Int, error)
	BlockEvents(ctx context.Context, startBlock uint64, endBlock uint64, eventTypes []string, contracts []string) ([]log.RawLog, error)
	CallContract(ctx context.Context, block *big.Int, sender, contract string, input []byte) ([]byte, error)
	Close()
}
