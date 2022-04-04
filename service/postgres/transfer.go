package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/events"
)

func (s *Store) UpsertTransferEvent(event events.Transfer) error {
	_, err := s.sqlBuilder.
		Insert(transferEventTableName).
		Columns(transferEventTableColumns...).
		Values(event.ID, event.Block, event.EventIndex, event.TransactionHash, event.CollectionID, event.TokenID, event.FromAddress, event.ToAddress, event.EmittedAt).
		Suffix(transferTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert transfer event: %w", err)
	}

	return nil
}
