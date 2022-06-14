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

func (t *TransferRepository) Upsert(transfers ...*events.Transfer) error {

	if len(transfers) == 0 {
		return nil
	}

	query := t.build.
		Insert("transfers").
		Columns(
			"id",
			"chain_id",
			"collection_address",
			"token_id",
			"block_number",
			"transaction_hash",
			"event_index",
			"sender_address",
			"receiver_address",
			"token_count",
			"emitted_at",
		)

	for _, transfer := range transfers {
		query = query.Values(
			transfer.ID,
			transfer.ChainID,
			transfer.CollectionAddress,
			transfer.TokenID,
			transfer.BlockNumber,
			transfer.TransactionHash,
			transfer.EventIndex,
			transfer.SenderAddress,
			transfer.ReceiverAddress,
			transfer.TokenCount,
			transfer.EmittedAt,
		)
		fmt.Println(transfer.TokenCount)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert transfer event: %w", err)
	}

	return nil
}
