package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/jobs"

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

func (m *MarketplaceRepository) Combinations(chainID uint64) ([]*jobs.Combination, error) {

	result, err := m.build.
		Select("networks.chain_id, networks_marketplaces.contract_address, events.event_hash, networks_marketplaces.start_height").
		From("networks, networks_marketplaces, marketplaces_standards, standards, standards_events, events").
		Where("networks.chain_id = ?", chainID).
		Where("networks_marketplaces.network_id = networks.id").
		Where("networks_marketplaces.marketplace_id = marketplaces_standards.marketplace_id").
		Where("marketplaces_standards.standard_id = standards.id").
		Where("standards.id = standards_events.standard_id").
		Where("standards_events.event_id = events.id").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	var combinations []*jobs.Combination
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var combination jobs.Combination
		err = result.Scan(
			&combination.ChainID,
			&combination.ContractAddress,
			&combination.EventHash,
			&combination.StartHeight,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		combinations = append(combinations, &combination)
	}

	return combinations, nil
}

func (m *MarketplaceRepository) RetrieveByAddress(chainID string, address string) (*graph.Marketplace, error) {

	result, err := m.build.
		Select("marketplaces.id", "marketplaces.name", "marketplaces.description", "marketplaces.website").
		From("marketplaces, chains_marketplaces").
		Where("LOWER(chains_marketplaces.address) = LOWER(?)", address).
		Where("chains_marketplaces.chain_id = ?", chainID).
		Where("chains_marketplaces.marketplace_id = marketplaces.id").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not get row: %w", result.Err())
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
