package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type TraitRepository struct {
	build squirrel.StatementBuilderType
}

func NewTraitRepository(db *sql.DB) *TraitRepository {

	cache := squirrel.NewStmtCache(db)
	s := TraitRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &s
}

func (t *TraitRepository) Insert(traits ...*graph.Trait) error {

	if len(traits) == 0 {
		return nil
	}

	query := t.build.
		Insert(TableTraits).
		Columns(ColumnsTraits...)

	for _, trait := range traits {
		query = query.Values(
			trait.ID,
			trait.NFTID,
			trait.Name,
			trait.Type,
			trait.Value,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert trait: %w", err)
	}

	return nil
}
