package pipeline

import (
	"github.com/NFT-com/indexer/models/graph"
)

type CollectionStore interface {
	One(chainID uint64, address string) (*graph.Collection, error)
}

type NFTStore interface {
	Insert(nft *graph.NFT) error
}

type OwnerStore interface {
	AddCount(nftID string, owner string, count int) error
}

type TraitStore interface {
	Insert(traits ...*graph.Trait) error
}
