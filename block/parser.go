package block

import (
	"github.com/NFT-com/indexer/events"
	"golang.org/x/net/context"
)

type Block struct {
	Hash string
}

func (b *Block) String() string {
	return b.Hash
}

type Parser interface {
	ParseBlock(ctx context.Context, block *Block) ([]events.Event, error)
}
