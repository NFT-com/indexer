package pipeline

import (
	"github.com/NFT-com/indexer/models/graph"
)

type CollectionStore interface {
	One(chainID string, address string) (*graph.Collection, error)
}

type NFTStore interface {
	Upsert(nft *graph.NFT) error
	ChangeOwner(chainID string, address string, tokenID string, owner string) error
}

type TraitStore interface {
	Upsert(traits ...*graph.Trait) error
}
