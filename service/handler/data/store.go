package data

import (
	"github.com/NFT-com/indexer/models/chain"
)

type Store interface {
	ChainStore
	CollectionStore
}

type ChainStore interface {
	CreateChain(chain chain.Chain) error
	Chains() ([]chain.Chain, error)
}

type CollectionStore interface {
	CreateCollection(collection chain.Collection) error
	Collections() ([]chain.Collection, error)
}
