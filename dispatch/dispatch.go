package dispatch

import (
	"context"

	"github.com/NFT-com/indexer/event"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, event *event.Event) error
}
