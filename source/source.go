package source

import (
	"context"

	"github.com/NFT-com/indexer/block"
)

type Source interface {
	Next(ctx context.Context) *block.Block
	Close() error
}
