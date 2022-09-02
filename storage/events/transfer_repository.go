package events

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/database"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/storage"
)

type TransferRepository struct {
	build   squirrel.StatementBuilderType
	retrier storage.Retrier
}

func NewTransferRepository(db *sql.DB, retrier storage.Retrier) *TransferRepository {

	cache := squirrel.NewStmtCache(db)
	t := TransferRepository{
		build:   squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
		retrier: retrier,
	}

	return &t
}

func (t *TransferRepository) Upsert(transfers ...*events.Transfer) error {

	for start := 0; start < len(transfers); start += database.BatchSize {

		end := start + database.BatchSize
		if end > len(transfers) {
			end = len(transfers)
		}

		batch := transfers[start:end]
		if len(batch) == 0 {
			continue
		}

		query := t.build.
			Insert("transfers").
			Columns(
				"id",
				"chain_id",
				"token_standard",
				"collection_address",
				"token_id",
				"block_number",
				"transaction_hash",
				"event_index",
				"sender_address",
				"receiver_address",
				"token_count",
				"emitted_at",
			).
			Suffix("ON CONFLICT DO NOTHING")

		for _, transfer := range batch {
			query = query.Values(
				transfer.ID,
				transfer.ChainID,
				transfer.TokenStandard,
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
		}

		err := t.retrier.Insert(query)
		if err != nil {
			return fmt.Errorf("could not upsert transfer batch (start: %d, end: %d): %w", start, end, err)
		}
	}

	return nil
}
