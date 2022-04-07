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
	BlockEvents(ctx context.Context, block, event, contract string) ([]log.RawLog, error)
	CallContract(ctx context.Context, block *big.Int, sender, contract string, input []byte) ([]byte, error)
	Close()
}
