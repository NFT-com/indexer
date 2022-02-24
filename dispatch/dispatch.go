package dispatch

import (
	"context"

	"github.com/NFT-com/indexer/block"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, event *block.Block) error
}
