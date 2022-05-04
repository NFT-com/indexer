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

func (b *BurnRepository) Upsert(burns ...*events.Burn) error {

	query := b.build.
		Insert(TableBurnEvents).
		Columns(ColumnsBurnEvents...).
		Suffix(ConflictBurnEvents)

	for _, burn := range burns {
		query = query.Values(
			burn.ID,
			burn.BlockNumber,
			burn.EventIndex,
			burn.CollectionAddress,
			burn.TransactionHash,
			burn.TokenID,
			burn.EmittedAt,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert burn event: %w",
			err)
	}

	return nil
}
