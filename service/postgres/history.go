package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/event"
	"github.com/lib/pq"
)

func (s *Store) InsertHistory(event event.Event) error {
	_, err := s.sqlBuilder.
		Insert(HistoryDBName).
		Columns(HistoryTableColumns...).
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
		if ok && pqErr.Code == UniqueViolation {
			return ErrAlreadyExists
		}

		return fmt.Errorf("failed to insert history: %v", err)
	}

	return nil
}
