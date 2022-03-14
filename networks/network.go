package networks

import (
	"context"
	"github.com/NFT-com/indexer/events"
)

const (
	EventTypeMint   = "mint"
	EventTypeUpdate = "update"
	EventTypeBurn   = "burn"
)

type Network interface {
	BlockEvents(ctx context.Context, block, event, contract string) ([]events.RawEvent, error)
	Close()
}
