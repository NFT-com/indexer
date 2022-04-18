package addition

import (
	"github.com/NFT-com/indexer/models/chain"
)

type Store interface {
	UpsertNFT(nft chain.NFT, collectionID string) error
	UpsertTrait(trait chain.Trait) error

	Chain(chainID string) (*chain.Chain, error)
	Collection(chainID, address, contractCollectionID string) (*chain.Collection, error)
}
