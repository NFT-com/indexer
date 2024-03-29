package pipeline

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
)

type NFTStore interface {
	Missing(touches ...*graph.NFT) ([]*graph.NFT, error)
	Touch(touches ...*graph.NFT) error
	Upsert(nft *graph.NFT) error
	Delete(deletions ...*graph.NFT) error
}

type TraitStore interface {
	Upsert(traits ...*graph.Trait) error
}

type OwnerStore interface {
	Upsert(transfers ...*events.Transfer) error
	Sanitize() error
}

type BoundaryStore interface {
	All() ([]*jobs.Boundary, error)
	Upsert(chainID uint64, addresses []string, events []string, height uint64, jobID string) error
}

type FailureStore interface {
	Parsing(parsing *jobs.Parsing, message string) error
	Addition(addition *jobs.Addition, message string) error
	Completion(completion *jobs.Completion, message string) error
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
	Update(sales ...*events.Sale) error
}
