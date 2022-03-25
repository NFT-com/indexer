package postgres

import (
	"fmt"
	"time"
)

func (s *Store) InsertNewNFT(chain, contract, id, owner string) error {
	_, err := s.sqlBuilder.
		Insert(nftsTableName).
		Columns(nftsTableColumns...).
		Values(id, chain, contract, owner).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert new nft: %w", err)
	}

	return nil
}

func (s *Store) UpdateNFT(chain, contract, id, owner string) error {
	_, err := s.sqlBuilder.
		Update(nftsTableName).
		Where("id = ? AND network_id = ? AND chain_id = ? AND contract = ?", id, chain, contract).
		Set("owner", owner).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update nft: %w", err)
	}

	return nil
}
