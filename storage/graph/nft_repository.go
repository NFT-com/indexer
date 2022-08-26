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

func (n *NFTRepository) Touch(touches ...*graph.NFT) error {

	if len(touches) == 0 {
		return nil
	}

	set := make(map[string]*graph.NFT, len(touches))
	for _, touch := range touches {
		set[touch.ID] = touch
	}

	touches = make([]*graph.NFT, 0, len(set))
	for _, touch := range set {
		touches = append(touches, touch)
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

	for _, touch := range touches {
		query = query.Values(
			touch.ID,
			touch.CollectionID,
			touch.TokenID,
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

func (n *NFTRepository) Delete(deletions ...*graph.NFT) error {

	if len(deletions) == 0 {
		return nil
	}

	set := make(map[string]*graph.NFT, len(deletions))
	for _, deletion := range deletions {
		set[deletion.ID] = deletion
	}

	deletions = make([]*graph.NFT, 0, len(set))
	for _, deletion := range set {
		deletions = append(deletions, deletion)
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
			"deleted",
			"deleted_at",
		).
		Suffix("ON CONFLICT (id) DO UPDATE SET " +
			"deleted = TRUE, " +
			"deleted_at = EXCLUDED.deleted_at " +
			"WHERE nfts.deleted = FALSE")

	for _, deletion := range deletions {
		query = query.Values(
			deletion.ID,
			deletion.CollectionID,
			deletion.TokenID,
			"",
			"",
			"",
			"",
			true,
			"NOW()",
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	return nil
}
