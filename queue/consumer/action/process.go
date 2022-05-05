package action

import (
	"fmt"

	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/models/chain"
)

func (d *Action) processNFT(actionType string, nft chain.NFT) error {

	chain, err := d.dataStore.Chain(nft.ChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}

	collection, err := d.dataStore.Collection(chain.ID, nft.Contract)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

	switch actionType {

	case log.Addition.String():

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

		d.log.Info().
			Str("collection", collection.ID).
			Str("nft", nft.ID).
			Str("name", nft.Name).
			Str("uri", nft.URI).
			Str("image", nft.Image).
			Int("traits", len(nft.Traits)).
			Msg("NFT details added")

	case log.OwnerChange.String():

		err = d.dataStore.UpdateNFTOwner(collection.ID, nft.ID, nft.Owner)
		if err != nil {
			return fmt.Errorf("could not update nft owner (nft %s): %w", nft.ID, err)
		}

		d.log.Info().
			Str("collection", collection.ID).
			Str("nft", nft.ID).
			Str("owner", nft.Owner).
			Msg("NFT owner updated")
	}

	return nil
}
