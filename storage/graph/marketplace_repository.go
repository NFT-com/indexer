package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type MarketplaceRepository struct {
	build squirrel.StatementBuilderType
}

func NewMarketplaceRepository(db *sql.DB) *MarketplaceRepository {

	cache := squirrel.NewStmtCache(db)
	m := MarketplaceRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &m
}

func (m *MarketplaceRepository) RetrieveByAddress(chainID string, address string) (*graph.Marketplace, error) {

	result, err := m.build.
		Select("marketplaces.id", "marketplaces.name", "marketplaces.description", "marketplaces.website").
		From("marketplaces, chains_marketplaces").
		Where("chains_marketplaces.address = ?", address).
		Where("chains_marketplaces.chain_id = ?", chainID).
		Where("chains_marketplaces.marketplace_id = marketplaces.id").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not get row: %w", err)
	}
	if !result.Next() {
		return nil, sql.ErrNoRows
	}

	var marketplace graph.Marketplace
	err = result.Scan(
		&marketplace.ID,
		&marketplace.Name,
		&marketplace.Description,
		&marketplace.Website,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	return &marketplace, nil
}
