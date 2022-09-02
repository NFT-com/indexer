package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/database"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/storage"
)

type TraitRepository struct {
	build   squirrel.StatementBuilderType
	retrier storage.Retrier
}

func NewTraitRepository(db *sql.DB, retrier storage.Retrier) *TraitRepository {

	cache := squirrel.NewStmtCache(db)
	s := TraitRepository{
		build:   squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
		retrier: retrier,
	}

	return &s
}

func (t *TraitRepository) Upsert(traits ...*graph.Trait) error {

	for start := 0; start < len(traits); start += database.BatchSize {

		end := start + database.BatchSize
		if end > len(traits) {
			end = len(traits)
		}

		batch := traits[start:end]
		if len(batch) == 0 {
			continue
		}

		query := t.build.
			Insert("traits").
			Columns(
				"id",
				"nft_id",
				"name",
				"type",
				"value",
			).
			Suffix("ON CONFLICT (id) DO UPDATE SET " +
				"name = EXCLUDED.name, " +
				"type = EXCLUDED.type, " +
				"value = EXCLUDED.value")

		for _, trait := range batch {
			query = query.Values(
				trait.ID,
				trait.NFTID,
				trait.Name,
				trait.Type,
				trait.Value,
			)
		}

		err := t.retrier.Insert(query)
		if err != nil {
			return fmt.Errorf("could not upsert trait batch (start: %d, end: %d): %w", start, end, err)
		}
	}

	return nil
}
