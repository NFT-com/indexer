package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) EventTypes(standardID string) ([]chain.EventType, error) {
	result, err := s.sqlBuilder.
		Select("id, name").
		From("event_types").
		Where("standard = ?", standardID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collections: %w", err)
	}
	defer result.Close()

	eventTypes := make([]chain.EventType, 0)
	for result.Next() && result.Err() == nil {
		var eventType chain.EventType
		err = result.Scan(
			&eventType.ID,
			&eventType.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan collection: %w", err)
		}

		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes, nil
}
