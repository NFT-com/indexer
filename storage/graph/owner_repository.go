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

func (n *OwnerRepository) Upsert(transfers ...*events.Transfer) error {

	if len(transfers) == 0 {
		return nil
	}

	query := n.build.
		Insert("owners").
		Columns(
			"owner",
			"nft_id",
			"event_id",
			"number",
		).
		Suffix("ON CONFLICT DO NOTHING")

	for _, transfer := range transfers {

		// skip transfers that have same sender & receiver
		if transfer.SenderAddress == transfer.ReceiverAddress {
			continue
		}

		query = query.Values(
			transfer.SenderAddress,
			transfer.NFTID(),
			transfer.EventID(),
			-int(transfer.TokenCount),
		)

		query = query.Values(
			transfer.ReceiverAddress,
			transfer.NFTID(),
			transfer.EventID(),
			int(transfer.TokenCount),
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}
