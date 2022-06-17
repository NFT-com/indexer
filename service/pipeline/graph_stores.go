package pipeline

import (
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
)

type CollectionStore interface {
	One(chainID uint64, address string) (*graph.Collection, error)
}

type NFTStore interface {
	Touch(nftID string, collectionID string, tokenID string) error
	Insert(nft *graph.NFT) error
}

type OwnerStore interface {
	Add(additions ...*jobs.Addition) error
	Change(modifications ...*jobs.Modification) error
}

type TraitStore interface {
	Insert(traits ...*graph.Trait) error
}
