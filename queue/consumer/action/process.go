package action

import (
	"fmt"

	"github.com/NFT-com/indexer/action"
	"github.com/NFT-com/indexer/models/chain"
)

func (d *Action) processNFT(actionType string, nft chain.NFT) error {
	chain, err := d.dataStore.Chain(nft.ChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}

	collection, err := d.dataStore.Collection(chain.ID, nft.Contract, nft.ContractCollectionID)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	switch actionType {

	case action.Addition:
		err = d.dataStore.UpsertNFT(nft, collection.ID)
		if err != nil {
			return fmt.Errorf("could not store nft: %w", err)
		}

		for _, trait := range nft.Traits {
			err = d.dataStore.UpsertTrait(trait)
			if err != nil {
				return fmt.Errorf("could not store trait: %w", err)
			}
		}

	case action.OwnerChange:
		err = d.dataStore.UpdateNFTOwner(collection.ID, nft.ID, nft.Owner)
		if err != nil {
			return fmt.Errorf("could not update nft owner (nft %s): %w", nft.ID, err)
		}
	}

	return nil
}
