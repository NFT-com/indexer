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

func (n *NFTRepository) Touch(dummies ...*graph.NFT) error {

	if len(dummies) == 0 {
		return nil
	}

	set := make(map[string]*graph.NFT, len(dummies))
	for _, dummy := range dummies {
		set[dummy.ID] = dummy
	}

	dummies = make([]*graph.NFT, 0, len(set))
	for _, dummy := range set {
		dummies = append(dummies, dummy)
	}

	query := n.build.
		Insert("nfts").
		Columns(
			"id",
			"collection_id",
			"token_id",
			"name",
			"uri",
			"image",
			"description",
			"updated_at",
		).
		Suffix("ON CONFLICT (id) DO UPDATE SET " +
			"updated_at = EXCLUDED.updated_at")

	for _, dummy := range dummies {
		query = query.Values(
			dummy.ID,
			dummy.CollectionID,
			dummy.TokenID,
			"",
			"",
			"",
			"",
			"NOW()",
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	return nil
}

func (n *NFTRepository) Upsert(nft *graph.NFT) error {

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
			"created_at",
		).
		Values(
			nft.ID,
			nft.CollectionID,
			nft.TokenID,
			nft.Name,
			nft.URI,
			nft.Image,
			nft.Description,
			"NOW()",
		).
		Suffix("ON CONFLICT (id) DO UPDATE SET " +
			"name = EXCLUDED.name, " +
			"uri = EXCLUDED.uri, " +
			"image = EXCLUDED.image, " +
			"description = EXCLUDED.description, " +
			"created_at = EXCLUDED.created_at").
		Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}
