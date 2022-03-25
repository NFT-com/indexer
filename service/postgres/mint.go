package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/events"
)

func (s *Store) UpsertMintEvent(event events.Mint) error {
	_, err := s.sqlBuilder.
		Insert(mintEventTableName).
		Columns(mintEventTableColumns...).
		Values(event.ID, event.Block, event.TransactionHash, event.CollectionID, event.Owner, event.EmittedAt).
		Suffix(mintTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert mint event: %w", err)
	}

	return nil
}

func (s *Store) UpsertMintEvents(events []events.Mint) error {
	query := s.sqlBuilder.
		Insert(mintEventTableName).
		Columns(mintEventTableColumns...).
		Suffix(mintTableOnConflictStatement)

	for _, event := range events {
		query = query.Values(event.ID, event.Block, event.TransactionHash, event.CollectionID, event.Owner, event.EmittedAt)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not insert mint events: %w", err)
	}

	return nil
}
