package networks

import (
	"context"

	"github.com/NFT-com/indexer/events"
)

type Network interface {
	BlockEvents(ctx context.Context, block, event, contract string) ([]events.RawEvent, error)
	Close()
}
