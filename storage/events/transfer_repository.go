package events

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/events"
)

type TransferRepository struct {
	build squirrel.StatementBuilderType
}

func NewTransferRepository(db *sql.DB) *TransferRepository {

	cache := squirrel.NewStmtCache(db)
	t := TransferRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &t
}

func (t *TransferRepository) UpsertTransferEvent(transfer *events.Transfer) error {

	_, err := t.build.
		Insert(TableTransferEvents).
		Columns(ColumnsTransferEvents...).
		Values(
			transfer.ID,
			transfer.Block,
			transfer.EventIndex,
			transfer.TransactionHash,
			transfer.CollectionID,
			transfer.TokenID,
			transfer.FromAddress,
			transfer.ToAddress,
			transfer.EmittedAt,
		).
		Suffix(ConflictTransferEvents).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert transfer event: %w",
			err)
	}

	return nil
}
