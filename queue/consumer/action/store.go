package action

import (
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
)

type Store interface {
	UpsertNFT(nft chain.NFT, collectionID string) error
	UpsertTrait(trait chain.Trait) error

	UpdateNFTOwner(collectionID, nft, owner string) error

	Chain(chainID string) (*chain.Chain, error)
	Collection(chainID, address, contractCollectionID string) (*chain.Collection, error)

	ActionJob(id string) (*jobs.Action, error)
	UpdateActionJobStatus(id string, status jobs.Status) error
}
