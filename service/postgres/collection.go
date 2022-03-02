package postgres

import (
	"database/sql"
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) CreateCollection(collection chain.Collection) error {
	_, err := s.sqlBuilder.
		Insert(collectionTableName).
		Columns(collectionTableColumns...).
		Values(
			collection.ID,
			collection.ChainID,
			collection.ContractCollectionID,
			collection.Address,
			collection.Name,
			collection.Description,
			collection.Symbol,
			collection.Slug,
			collection.URI,
			collection.ImageURL,
			collection.Website,
		).
		Exec()

	if err != nil {
		return fmt.Errorf("could not create chain: %w", err)
	}

	return nil
}

func (s *Store) Collections() ([]chain.Collection, error) {
	result, err := s.sqlBuilder.
		Select(collectionTableColumns...).
		From(collectionTableName).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve collection list: %w", err)
	}

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
			&collection.URI,
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
	query := s.sqlBuilder.
		Select(collectionTableColumns...).
		From(collectionTableName).
		Where("chain_id = ?", chainID).
		Where("address = ?", address)

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
		&collection.URI,
		&collection.ImageURL,
		&collection.Website,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan collection: %w", err)
	}

	collection.ContractCollectionID = ccID.String
	return &collection, nil
}
