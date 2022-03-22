package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"github.com/NFT-com/indexer/events"
)

func (s *Store) InsertRawEvent(event events.RawEvent) error {
	columns := make([]string, 0, len(eventsTableColumns))
	copy(columns, eventsTableColumns)
	values := []interface{}{
		event.ID,
		event.ChainID,
		event.NetworkID,
		event.BlockNumber,
		event.BlockHash,
		event.Address,
		event.TransactionHash,
		event.EventType,
	}

	if len(event.IndexData) > 0 {
		indexData, err := json.Marshal(event.IndexData)
		if err != nil {
			return err
		}

		columns = append(columns, eventsTableIndexedDataColumn)
		values = append(values, indexData)
	}

	if len(event.Data) > 0 {
		columns = append(columns, eventsTableDataColumn)
		values = append(values, event.Data)
	}

	_, err := s.sqlBuilder.
		Insert(eventsTableName).
		Columns(columns...).
		Values(values...).
		Exec()
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == uniqueViolation {
			return errAlreadyExists
		}

		return fmt.Errorf("could not insert events: %w", err)
	}

	return nil
}
