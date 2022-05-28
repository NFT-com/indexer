package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type NFTRepository struct {
	build squirrel.StatementBuilderType
}

func NewNFTRepository(db *sql.DB) *NFTRepository {

	cache := squirrel.NewStmtCache(db)
	n := NFTRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &n
}

func (n *NFTRepository) Touch(nftID string) error {

	_, err := n.build.
		Insert("nfts").
		Columns("id").
		Values(nftID).
		Suffix("ON CONFLICT (id) DO NOTHING").
		Exec()

	if err != nil {
		return fmt.Errorf("could not execute statement: %w", err)
	}
	return nil
}

func (n *NFTRepository) Insert(nft *graph.NFT) error {

	_, err := n.build.
		Insert("nfts").
		Columns(
			"id",
			"collection_id",
			"token_id",
			"name",
			"uri",
			"image",
			"description",
		).
		Values(
			nft.ID,
			nft.CollectionID,
			nft.TokenID,
			nft.Name,
			nft.URI,
			nft.Image,
			nft.Description,
		).
		Suffix("ON CONFLICT (id) DO UPDATE SET " +
			"collection_id = EXCLUDED.collection_id, " +
			"token_id = EXCLUDED.token_id, " +
			"name = EXCLUDED.name, " +
			"uri = EXCLUDED.uri, " +
			"image = EXCLUDED.image, " +
			"description = EXCLUDED.description").
		Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}
