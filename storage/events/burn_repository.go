package events

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/events"
)

type BurnRepository struct {
	build squirrel.StatementBuilderType
}

func NewBurnRepository(db *sql.DB) *BurnRepository {

	cache := squirrel.NewStmtCache(db)
	b := BurnRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &b
}

func (b *BurnRepository) Upsert(event *events.Burn) error {

	_, err := b.build.
		Insert(TableBurnEvents).
		Columns(ColumnsBurnEvents...).
		Values(
			event.ID,
			event.Block,
			event.EventIndex,
			event.TransactionHash,
			event.CollectionID,
			event.TokenID,
			event.EmittedAt,
		).
		Suffix(ConflictBurnEvents).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert burn event: %w",
			err)
	}

	return nil
}
