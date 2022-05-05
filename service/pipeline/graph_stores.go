package pipeline

import (
	"github.com/NFT-com/indexer/models/graph"
)

type NetworkStore interface {
	Retrieve(chainID string) (*graph.Network, error)
}

type CollectionStore interface {
	RetrieveByAddress(chainID string, address string, collectionID string) (*graph.Collection, error)
}

type MarketplaceStore interface {
	RetrieveByAddress(chainID string, address string) (*graph.Marketplace, error)
}

type NFTStore interface {
	Upsert(nft *graph.NFT) error
	ChangeOwner(nftID string, owner string) error
}

type TraitStore interface {
	Upsert(traits ...*graph.Trait) error
}
