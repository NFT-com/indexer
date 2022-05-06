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

func (n *NFTRepository) Insert(nft *graph.NFT) error {

	_, err := n.build.
		Insert(TableNFTs).
		Columns(ColumnsNFTs...).
		Values(
			nft.ID,
			nft.CollectionID,
			nft.TokenID,
			nft.Name,
			nft.URI,
			nft.Image,
			nft.Description,
			nft.Owner,
		).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}

func (n *NFTRepository) ChangeOwner(collectionID string, tokenID string, owner string) error {

	_, err := n.build.
		Update("nfts").
		Set("owner", owner).
		Set("updated_at", "NOW()").
		Where("collection_id = ?", collectionID).
		Where("token_id = ?", tokenID).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}
