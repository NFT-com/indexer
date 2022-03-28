package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/marketplace"
)

func (s *Store) Marketplace(chainID, address string) (*marketplace.Marketplace, error) {
	result, err := s.sqlBuilder.
		Select("marketplaces.id", "marketplaces.name", "marketplaces.description", "marketplaces.website").
		From("marketplaces, chains_marketplaces").
		Where("chains_marketplaces.address = ?", address).
		Where("chains_marketplaces.chain_id = ?", chainID).
		Where("chains_marketplaces.marketplace_id = marketplaces.id").
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve marketplace: %w", errResourceNotFound)
	}

	var marketplace marketplace.Marketplace
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
