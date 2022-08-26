package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

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

	if len(transfers) == 0 {
		return nil
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

	for _, transfer := range transfers {

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
			return fmt.Errorf("could not delete row: %w", err)
		}
	}

	return nil
}
