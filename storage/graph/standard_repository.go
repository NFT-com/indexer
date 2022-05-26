package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type StandardRepository struct {
	build squirrel.StatementBuilderType
}

func NewStandardRepository(db *sql.DB) *StandardRepository {

	cache := squirrel.NewStmtCache(db)
	s := StandardRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &s
}

func (s *StandardRepository) List() ([]*graph.Standard, error) {

	result, err := s.build.
		Select("*").
		From("standards").
		OrderBy("id ASC").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not list collections: %w", err)
	}
	defer result.Close()

	var standards []*graph.Standard
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var standard graph.Standard
		err = result.Scan(
			&standard.ID,
			&standard.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		standards = append(standards, &standard)
	}

	return standards, nil
}
