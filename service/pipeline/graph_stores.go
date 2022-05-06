package pipeline

import (
	"github.com/NFT-com/indexer/models/graph"
)

type CollectionStore interface {
	One(chainID uint64, address string) (*graph.Collection, error)
}

type NFTStore interface {
	Insert(nft *graph.NFT) error
	ChangeOwner(collectionID string, tokenID string, owner string) error
}

type TraitStore interface {
	Insert(traits ...*graph.Trait) error
}
