package postgres

import (
	"fmt"

	"github.com/lib/pq"

	"github.com/NFT-com/indexer/event"
)

func (s *Store) InsertHistory(event event.Event) error {
	_, err := s.sqlBuilder.
		Insert(historyTableName).
		Columns(historyTableColumns...).
		Values(
			event.ID,
			event.ChainID,
			event.NetworkID,
			event.Type,
			event.Contract,
			event.NftID,
			event.FromAddress,
			event.ToAddress,
			event.Price,
			event.EmittedAt,
		).
		Exec()
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == uniqueViolation {
			return errAlreadyExists
		}

		return fmt.Errorf("failed to insert history: %v", err)
	}

	return nil
}
