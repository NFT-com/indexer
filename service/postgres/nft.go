package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) UpsertNFT(nft chain.NFT, collectionID string) error {
	_, err := s.sqlBuilder.
		Insert(nftTableName).
		Columns(nftTableColumns...).
		Values(nft.ID, nft.TokenID, collectionID, nft.Name, nft.Image, nft.Description, nft.Owner).
		Suffix(nftTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}
