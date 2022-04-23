package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/events"
)

func (s *Store) UpsertBurnEvent(event events.Burn) error {

	_, err := s.build.
		Insert(burnEventTableName).
		Columns(burnEventTableColumns...).
		Values(event.ID, event.Block, event.EventIndex, event.TransactionHash, event.CollectionID, event.TokenID, event.EmittedAt).
		Suffix(burnTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert burn event: %w", err)
	}

	return nil
}
