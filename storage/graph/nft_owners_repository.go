package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type NFTOwnerRepository struct {
	build squirrel.StatementBuilderType
}

func NewNFTOwnerRepository(db *sql.DB) *NFTOwnerRepository {

	cache := squirrel.NewStmtCache(db)
	n := NFTOwnerRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &n
}

func (n *NFTOwnerRepository) Upsert(nfts ...*graph.NFT) error {

	query := n.build.
		Insert(TableNFTOwners).
		Columns(ColumnsNFTOwners...).
		Suffix(ConflictNFTOwners)

	for _, nft := range nfts {
		query = query.Values(
			nft.ID,
			nft.Owner,
			nft.Number,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}
