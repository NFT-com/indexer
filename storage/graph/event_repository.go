package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type EventRepository struct {
	build squirrel.StatementBuilderType
}

func NewEventRepository(db *sql.DB) *EventRepository {

	cache := squirrel.NewStmtCache(db)
	c := EventRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &c
}

func (e *EventRepository) ListForStandard(standardID string) ([]*graph.Event, error) {

	result, err := e.build.
		Select("events.id", "events.event_hash", "events.name").
		From("events, standards_events").
		Where("events.id = standards_events.event_id").
		Where("standards_events.standard_id = ?", standardID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	events := make([]*graph.Event, 0)
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var event graph.Event
		err = result.Scan(
			&event.ID,
			&event.EventHash,
			&event.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		events = append(events, &event)
	}

	return events, nil
}
