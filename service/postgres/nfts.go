package postgres

import (
	"fmt"
	"time"
)

func (s *Store) InsertNewNFT(network, chain, contract, id, owner string) error {
	_, err := s.sqlBuilder.
		Insert(NFTsDBName).
		Columns(NFTsTableColumns...).
		Values(id, network, chain, contract, owner).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to insert new nft: %v", err)
	}

	return nil
}

func (s *Store) UpdateNFT(network, chain, contract, id, owner string) error {
	_, err := s.sqlBuilder.
		Update(NFTsDBName).
		Where("id = ? AND network_id = ? AND chain_id = ? AND contract = ?", id, network, chain, contract).
		Set("owner", owner).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to update nft: %v", err)
	}

	return nil
}
