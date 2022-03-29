package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) Marketplace(chainID, address string) (*chain.Marketplace, error) {
	result, err := s.sqlBuilder.
		Select("marketplaces.id", "marketplaces.name", "marketplaces.description", "marketplaces.website").
		From("marketplaces, chains_marketplaces").
		Where("chains_marketplaces.address = ?", address).
		Where("chains_marketplaces.chain_id = ?", chainID).
		Where("chains_marketplaces.marketplace_id = marketplaces.id").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not query marketplace: %w", err)
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve marketplace: %w", errResourceNotFound)
	}

	var marketplace chain.Marketplace
	err = result.Scan(
		&marketplace.ID,
		&marketplace.Name,
		&marketplace.Description,
		&marketplace.Website,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan marketplace: %w", err)
	}

	return &marketplace, nil
}
