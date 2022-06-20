package pipeline

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type NFTStore interface {
	Touch(modifications ...*jobs.Modification) error
	Insert(nft *graph.NFT) error
}

type TraitStore interface {
	Insert(traits ...*graph.Trait) error
}

type OwnerStore interface {
	Add(additions ...*results.Addition) error
	Change(modifications ...*jobs.Modification) error
}

type BoundaryStore interface {
	ForCombination(chainID uint64, address string, event string) (uint64, error)
	Upsert(chainID uint64, addresses []string, events []string, height uint64, jobID string) error
}

type CollectionStore interface {
	One(chainID uint64, address string) (*graph.Collection, error)
	Combinations(chainID uint64) ([]*jobs.Combination, error)
}

type MarketplaceStore interface {
	Combinations(chainID uint64) ([]*jobs.Combination, error)
}

type TransferStore interface {
	Upsert(transfers ...*events.Transfer) error
}

type SaleStore interface {
	Upsert(sales ...*events.Sale) error
}
