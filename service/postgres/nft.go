package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) UpsertNFT(nft chain.NFT, collectionID string) error {

	_, err := s.build.
		Insert(nftTableName).
		Columns(nftTableColumns...).
		Values(
			nft.ID,
			nft.TokenID,
			collectionID,
			nft.Name,
			nft.URI,
			nft.Image,
			nft.Description,
			nft.Owner,
		).
		Suffix(nftTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}

func (s *Store) UpdateNFTOwner(collectionID, nft, owner string) error {

	_, err := s.build.
		Update(nftTableName).
		Set("owner", owner).
		Where("token_id = ?", nft).
		Where("collection = ?", collectionID).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}
