package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/NFT-com/indexer/events"
)

func (s *Store) InsertRawEvent(event events.RawEvent) error {
	columns := make([]string, len(EventsTableColumns))
	copy(columns, EventsTableColumns)
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

		columns = append(columns, EventsTableIndexedDataColumn)
		values = append(values, indexData)
	}

	if len(event.Data) > 0 {
		columns = append(columns, EventsTableDataColumn)
		values = append(values, event.Data)
	}

	_, err := s.sqlBuilder.
		Insert(EventsDBName).
		Columns(columns...).
		Values(values...).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to insert events: %v", err)
	}

	return nil
}
