package addition

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (d *Addition) processNFT(nft chain.NFT) error {
	chain, err := d.dataStore.Chain(nft.ChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}

	collection, err := d.dataStore.Collection(chain.ID, nft.Contract, nft.ContractCollectionID)
	if err != nil {
		return fmt.Errorf("could not get collection: %w", err)
	}

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
		Str("nft", nft.ID).
		Str("chain", nft.ChainID).
		Str("name", nft.Name).
		Msg("processed NFT")

	return nil
}
