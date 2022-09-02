package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/NFT-com/indexer/models/database"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/storage"
)

type NFTRepository struct {
	build   squirrel.StatementBuilderType
	retrier storage.Retrier
}

func NewNFTRepository(db *sql.DB, retrier storage.Retrier) *NFTRepository {

	cache := squirrel.NewStmtCache(db)
	n := NFTRepository{
		build:   squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
		retrier: retrier,
	}

	return &n
}

func (n *NFTRepository) Missing(touches ...*graph.NFT) ([]*graph.NFT, error) {

	set := make(map[string]struct{}, len(touches))
	for _, touch := range touches {
		set[touch.ID] = struct{}{}
	}

	nftIDs := make([]string, 0, len(set))
	for nftID := range set {
		nftIDs = append(nftIDs, nftID)
	}

	result, err := n.build.
		Select("id").
		From("nfts").
		Where("id = ANY(?)", pq.Array(nftIDs)).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	existing := make(map[string]struct{})
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var nftID string
		err = result.Scan(&nftID)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		existing[nftID] = struct{}{}
	}

	var filtered []*graph.NFT
	for _, touch := range touches {
		_, ok := existing[touch.ID]
		if ok {
			continue
		}
		filtered = append(filtered, touch)
	}

	return filtered, nil
}

func (n *NFTRepository) Touch(touches ...*graph.NFT) error {

	set := make(map[string]*graph.NFT, len(touches))
	for _, touch := range touches {
		set[touch.ID] = touch
	}

	touches = make([]*graph.NFT, 0, len(set))
	for _, touch := range set {
		touches = append(touches, touch)
	}

	for start := 0; start < len(touches); start += database.BatchSize {

		end := start + database.BatchSize
		if end > len(touches) {
			end = len(touches)
		}

		batch := touches[start:end]
		if len(batch) == 0 {
			continue
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

		for _, touch := range batch {
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

		err := n.retrier.Insert(query)
		if err != nil {
			return fmt.Errorf("could not execute touches (start: %d, end: %d): %w", start, end, err)
		}
	}

	return nil
}

func (n *NFTRepository) Upsert(nft *graph.NFT) error {

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
			"created_at = EXCLUDED.created_at")

	err := n.retrier.Insert(query)
	if err != nil {
		return fmt.Errorf("could not execute upsert: %w", err)
	}

	return nil
}

func (n *NFTRepository) Delete(deletions ...*graph.NFT) error {

	set := make(map[string]*graph.NFT, len(deletions))
	for _, deletion := range deletions {
		set[deletion.ID] = deletion
	}

	deletions = make([]*graph.NFT, 0, len(set))
	for _, deletion := range set {
		deletions = append(deletions, deletion)
	}

	for start := 0; start < len(deletions); start += database.BatchSize {

		end := start + database.BatchSize
		if end > len(deletions) {
			end = len(deletions)
		}

		batch := deletions[start:end]
		if len(batch) == 0 {
			continue
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

		for _, deletion := range batch {
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

		err := n.retrier.Insert(query)
		if err != nil {
			return fmt.Errorf("could not execute deletions (start: %d, end: %d): %w", start, end, err)
		}
	}

	return nil
}
