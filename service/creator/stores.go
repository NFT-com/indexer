package creator

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type BoundaryStore interface {
	Last(chainID uint64, address string, event string) (uint64, error)
	Upsert(chainID uint64, addresses []string, events []string, height uint64) error
}

type CollectionStore interface {
	Combinations(chainID uint64) ([]*jobs.Combination, error)
}

type MarketplaceStore interface {
	Combinations(chainID uint64) ([]*jobs.Combination, error)
}
