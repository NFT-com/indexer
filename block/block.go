package block

import (
	"golang.org/x/net/context"

	"github.com/NFT-com/indexer/event"
)

type Block string

func (b *Block) String() string {
	return string(*b)
}

type Parser interface {
	Parse(ctx context.Context, block *Block) ([]*event.Event, error)
}
