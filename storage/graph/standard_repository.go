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

func (s *StandardRepository) Find(collectionID string) ([]*graph.Standard, error) {

	result, err := s.build.
		Select("standards.id, standards.name").
		From("standards_collections, standards").
		Where("standards_collections.collection = ?", collectionID).
		Where("standards_collections.standard = standards.id").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collections: %w", err)
	}
	defer result.Close()

	var standards []*graph.Standard
	for result.Next() && result.Err() == nil {
		var standard graph.Standard
		err = result.Scan(
			&standard.ID,
			&standard.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve standards list: %w", err)
		}

		standards = append(standards, &standard)
	}

	return standards, nil
}
