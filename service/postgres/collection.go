package postgres

import (
	"fmt"
	"strings"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) Collections(chainID string) ([]chain.Collection, error) {

	result, err := s.build.
		Select(collectionTableColumns...).
		From(collectionTableName).
		Where("chain_id = ?", chainID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collections: %w", err)
	}
	defer result.Close()

	collections := make([]chain.Collection, 0)
	for result.Next() && result.Err() == nil {
		var (
			collection chain.Collection
		)

		err = result.Scan(
			&collection.ID,
			&collection.ChainID,
			&collection.Address,
			&collection.Name,
			&collection.Description,
			&collection.Symbol,
			&collection.Slug,
			&collection.ImageURL,
			&collection.Website,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan collection: %w", err)
		}

		collections = append(collections, collection)
	}

	return collections, nil
}

func (s *Store) Collection(chainID, address string) (*chain.Collection, error) {
	query := s.build.
		Select(collectionTableColumns...).
		From(collectionTableName).
		Where("chain_id = ?", chainID).
		Where("address = ?", strings.ToLower(address)) // FIXME: All addresses should be lowercased in all similar queries.

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collection: %w", err)
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve collection: %w", ErrResourceNotFound)
	}

	var (
		collection chain.Collection
	)

	err = result.Scan(
		&collection.ID,
		&collection.ChainID,
		&collection.Address,
		&collection.Name,
		&collection.Description,
		&collection.Symbol,
		&collection.Slug,
		&collection.ImageURL,
		&collection.Website,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan collection: %w", err)
	}

	return &collection, nil
}
