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

func (n *NFTRepository) UpsertNFT(nft *graph.NFT, collectionID string) error {

	_, err := n.build.
		Insert(TableNFTs).
		Columns(ColumnsNFTs...).
		Values(
			nft.ID,
			nft.TokenID,
			collectionID,
			nft.Name,
			nft.URI,
			nft.Image,
			nft.Description,
			nft.Owner,
		).
		Suffix(ConflictNFTs).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}

func (n *NFTRepository) UpdateNFT(nft *graph.NFT) error {

	_, err := n.build.
		Update(TableNFTs).
		Set("owner", nft.Owner).
		Where("id = ?", nft.ID).
		// Where("collection = ?", collectionID). // FIXME: doesn't seem to exist as a column?
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert nft: %w", err)
	}

	return nil
}
