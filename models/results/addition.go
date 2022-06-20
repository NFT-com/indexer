package results

import (
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
)

type Addition struct {
	Job      *jobs.Addition `json:"job"`
	NFT      *graph.NFT     `json:"nft"`
	Traits   []*graph.Trait `json:"traits"`
	Requests uint           `json:"requests"`
}
