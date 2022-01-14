package subscriber

import (
	"context"

	"github.com/NFT-com/indexer/events"
)

type Subscriber interface {
	Subscribe(ctx context.Context, events chan events.Event) error
	Status(ctx context.Context) error
	Close() error
}
