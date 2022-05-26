package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
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

func (n *OwnerRepository) AddCount(nftID string, owner string, count uint) error {
	query := n.build.
		Insert("owners").
		Columns("nft_id", "owner", "number").
		Values(
			nftID,
			owner,
			count,
		).
		Suffix("ON CONFLICT (nft_id, owner) DO UPDATE SET number = (owners.number + EXCLUDED.number)")

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}

func (n *OwnerRepository) MoveCount(nftID string, from string, to string, count uint) error {
	_, err := n.build.
		Select(
			fmt.Sprintf("transfer_tokens('%s', '%s', '%s', %d)", from, to, nftID, count),
		).
		Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}
