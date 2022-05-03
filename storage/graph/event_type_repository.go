package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type EventTypeRepository struct {
	build squirrel.StatementBuilderType
}

func NewEventTypeRepository(db *sql.DB) *EventTypeRepository {

	cache := squirrel.NewStmtCache(db)
	c := EventTypeRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &c
}

func (e *EventTypeRepository) EventTypes(wheres ...string) ([]*graph.EventType, error) {

	statement := e.build.
		Select(ColumnsEventTypes...).
		From(TableEventTypes)

	for _, where := range wheres {
		statement = statement.Where(where)
	}

	result, err := statement.Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	eventTypes := make([]*graph.EventType, 0)
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var eventType graph.EventType
		err = result.Scan(
			&eventType.ID,
			&eventType.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		eventTypes = append(eventTypes, &eventType)
	}

	return eventTypes, nil
}
