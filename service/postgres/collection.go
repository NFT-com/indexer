package postgres

import (
	"database/sql"
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
			ccID       sql.NullString
		)

		err = result.Scan(
			&collection.ID,
			&collection.ChainID,
			&ccID,
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

		collection.ContractCollectionID = ccID.String

		collections = append(collections, collection)
	}

	return collections, nil
}

func (s *Store) Collection(chainID, address, contractCollectionID string) (*chain.Collection, error) {
	query := s.build.
		Select(collectionTableColumns...).
		From(collectionTableName).
		Where("chain_id = ?", chainID).
		Where("address = ?", strings.ToLower(address)) // FIXME: All addresses should be lowercased in all similar queries.

	if contractCollectionID != "" {
		query = query.Where("contract_collection_id = ?", contractCollectionID)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collection: %w", err)
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve collection: %w", errResourceNotFound)
	}

	var (
		collection chain.Collection
		ccID       sql.NullString
	)

	err = result.Scan(
		&collection.ID,
		&collection.ChainID,
		&ccID,
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

	collection.ContractCollectionID = ccID.String
	return &collection, nil
}
