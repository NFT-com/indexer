package pipeline

import (
	"github.com/NFT-com/indexer/models/graph"
)

type ChainStore interface {
	Retrieve(chainID string) (*graph.Chain, error)
}

type CollectionStore interface {
	RetrieveByAddress(chainID string, address string, collectionID string) (*graph.Collection, error)
}

type MarketplaceStore interface {
	RetrieveByAddress(chainID string, address string) (*graph.Marketplace, error)
}
