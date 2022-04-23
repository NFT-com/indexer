package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) UpsertTrait(trait chain.Trait) error {

	_, err := s.build.
		Insert(traitsTableName).
		Columns(traitsTableColumns...).
		Values(trait.ID, trait.Name, trait.Value, trait.NftID).
		Suffix(traitsTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert trait: %w", err)
	}

	return nil
}
