package results

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/graph"
)

type Action struct {
	NFT    *graph.NFT     `json:"nft"`
	Traits []*graph.Trait `json:"traits"`
	Sale   *events.Sale   `json:"sale"`
}
