package networks

import (
	"context"
)

type Network interface {
	BlockEvents(ctx context.Context, block, event, contract string) ([]event.RawEvent, error)
	Close()
}
