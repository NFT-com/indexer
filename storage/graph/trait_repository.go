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

func (t *TraitRepository) UpsertTrait(trait *graph.Trait) error {

	_, err := t.build.
		Insert(TableTraits).
		Columns(ColumnsTraits...).
		Values(trait.ID, trait.Name, trait.Value, trait.NftID).
		Suffix(ConflictTraits).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert trait: %w", err)
	}

	return nil
}
