package networks

import (
	"context"

	"github.com/NFT-com/indexer/log"
)

type Network interface {
	BlockEvents(ctx context.Context, block, event, contract string) ([]log.RawLog, error)
	Close()
}
