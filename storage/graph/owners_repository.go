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

func (n *OwnerRepository) AddCount(nftID string, owner string, count int) error {

	_, err := n.build.
		Insert("owners").
		Columns("nft_id", "owner", "number").
		Values(nftID, owner, count).
		Suffix("ON CONFLICT (nft_id, owner) DO UPDATE SET number = (owners.number + EXCLUDED.number)").
		Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}
