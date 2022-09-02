package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/database"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/storage"
)

type OwnerRepository struct {
	build   squirrel.StatementBuilderType
	retrier storage.Retrier
}

func NewOwnerRepository(db *sql.DB, retrier storage.Retrier) *OwnerRepository {

	cache := squirrel.NewStmtCache(db)
	n := OwnerRepository{
		build:   squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
		retrier: retrier,
	}

	return &n
}

func (o *OwnerRepository) Upsert(transfers ...*events.Transfer) error {

	for start := 0; start < len(transfers); start += database.BatchSize {

		end := start + database.BatchSize
		if end > len(transfers) {
			end = len(transfers)
		}

		batch := transfers[start:end]
		if len(batch) == 0 {
			continue
		}

		query := o.build.
			Insert("owners").
			Columns(
				"owner",
				"nft_id",
				"event_id",
				"number",
			).
			Suffix("ON CONFLICT DO NOTHING")

		for _, transfer := range batch {

			query = query.Values(
				transfer.SenderAddress,
				transfer.NFTID(),
				transfer.EventID(),
				fmt.Sprintf("-%s", transfer.TokenCount),
			)

			query = query.Values(
				transfer.ReceiverAddress,
				transfer.NFTID(),
				transfer.EventID(),
				transfer.TokenCount,
			)
		}

		err := o.retrier.Insert(query)
		if err != nil {
			return fmt.Errorf("could not execute query: %w", err)
		}
	}

	return nil
}

func (o *OwnerRepository) Sanitize() error {

	result, err := o.build.
		Select("owner", "nft_id").
		From("owners").
		GroupBy("owner, nft_id").
		Having("SUM(number) = 0").
		Query()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	for result.Next() {

		if result.Err() != nil {
			return fmt.Errorf("could not get next row: %w", result.Err())
		}

		var owner, nftID string
		err = result.Scan(
			&owner,
			&nftID,
		)
		if err != nil {
			return fmt.Errorf("could not scan next row: %w", err)
		}

		query := o.build.
			Delete("owners").
			Where("owner = ?", owner).
			Where("nft_id = ?", nftID)

		err := o.retrier.Delete(query)
		if err != nil {
			return fmt.Errorf("could not delete rows: %w", err)
		}
	}

	return nil
}
