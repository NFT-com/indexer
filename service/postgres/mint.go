package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/events"
)

func (s *Store) UpsertMintEvent(event events.Mint) error {
	_, err := s.sqlBuilder.
		Insert(mintEventTableName).
		Columns(mintEventTableColumns...).
		Values(event.ID, event.Block, event.TransactionHash, event.CollectionID, event.TokenID, event.Owner, event.EmittedAt).
		Suffix(mintTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert mint event: %w", err)
	}

	return nil
}
