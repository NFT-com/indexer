package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/database"
	"github.com/NFT-com/indexer/models/events"
)

type OwnerRepository struct {
	build squirrel.StatementBuilderType
}

func NewOwnerRepository(db *sql.DB) *OwnerRepository {

	cache := squirrel.NewStmtCache(db)
	n := OwnerRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &n
}

func (o *OwnerRepository) Upsert(transfers ...*events.Transfer) error {

	for start := 0; start > len(transfers); start += database.BatchSize {

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

		_, err := query.Exec()
		if err != nil {
			return fmt.Errorf("could not execute query: %w", err)
		}
	}

	return nil
}

func (o *OwnerRepository) Sanitize() error {

	result, err := o.build.
		Select(
			"a.owner",
			"a.nft_id",
			"a.event_id",
			"b.event_id",
		).
		From("owners as a, owners as b").
		Where("a.owner = b.owner").
		Where("a.nft_id = b.nft_id").
		Where("a.number = -b.number").
		Query()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	for result.Next() {

		if result.Err() != nil {
			return fmt.Errorf("could not get next row: %w", result.Err())
		}

		var owner, nftID, eventID1, eventID2 string
		err = result.Scan(
			&owner,
			&nftID,
			&eventID1,
			&eventID2,
		)
		if err != nil {
			return fmt.Errorf("could not scan next row: %w", err)
		}

		_, err = o.build.
			Delete("owners").
			Where("owner = ?", owner).
			Where("nft_id = ?", nftID).
			Where("(event_id = ? OR event_id = ?)", eventID1, eventID2).
			Exec()
		if err != nil {
			return fmt.Errorf("could not delete rows: %w", err)
		}
	}

	return nil
}
