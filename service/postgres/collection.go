package postgres

import (
	"database/sql"
	"fmt"

	"github.com/NFT-com/indexer/collection"
)

func (s *Store) Collection(chainID, address, contractCollectionID string) (*collection.Collection, error) {
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
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve collection: %w", errResourceNotFound)
	}

	var (
		collection collection.Collection
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
		&collection.Standard,
		&collection.URI,
		&collection.ImageURL,
	)

	collection.ContractCollectionID = ccID.String
	if err != nil {
		return nil, fmt.Errorf("could not retrieve collection: %w", err)
	}

	return &collection, nil
}
