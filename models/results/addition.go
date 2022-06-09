package results

import (
	"github.com/NFT-com/indexer/models/graph"
)

type Addition struct {
	NFT      *graph.NFT     `json:"nft"`
	Traits   []*graph.Trait `json:"traits"`
	Requests uint           `json:"requests"`
}
