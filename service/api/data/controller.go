package data

import (
	"github.com/NFT-com/indexer/models/chain"
)

type DataController interface {
	ChainController
	CollectionController
}

type ChainController interface {
	CreateChain(chain chain.Chain) (*chain.Chain, error)
	ListChains() ([]chain.Chain, error)
}

type CollectionController interface {
	CreateCollection(collection chain.Collection) (*chain.Collection, error)
	ListCollections() ([]chain.Collection, error)
}
