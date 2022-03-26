package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/events"
)

func (s *Store) UpsertBurnEvent(event events.Burn) error {
	_, err := s.sqlBuilder.
		Insert(burnEventTableName).
		Columns(burnEventTableColumns...).
		Values(event.ID, event.Block, event.TransactionHash, event.CollectionID, event.EmittedAt).
		Suffix(burnTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert burn event: %w", err)
	}

	return nil
}
