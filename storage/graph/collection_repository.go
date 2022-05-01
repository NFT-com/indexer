package graph

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/storage/filters"
)

type CollectionRepository struct {
	build squirrel.StatementBuilderType
}

func NewCollectionRepository(db *sql.DB) *CollectionRepository {

	cache := squirrel.NewStmtCache(db)
	c := CollectionRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &c
}

func (c *CollectionRepository) RetrieveByAddress(chainID string, address string, contractCollectionID string) (*graph.Collection, error) {

	query := c.build.
		Select(ColumnsCollections...).
		From(TableCollections).
		Where("chain_id = ?", chainID).
		Where("address = ?", strings.ToLower(address))
	if contractCollectionID != "" {
		query = query.Where("contract_collection_id = ?", contractCollectionID)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collection: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not get row: %w", err)
	}
	if !result.Next() {
		return nil, sql.ErrNoRows
	}

	var collection graph.Collection
	var ccID sql.NullString
	err = result.Scan(
		&collection.ID,
		// &collection.CollectionID, // TODO: why did this field go away?
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
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	collection.ContractCollectionID = ccID.String
	return &collection, nil
}

func (c *CollectionRepository) Find(wheres ...filters.Where) ([]*graph.Collection, error) {

	statement := c.build.
		Select(ColumnsCollections...).
		From(TableCollections)

	for _, where := range wheres {
		statement = statement.Where(where())
	}

	result, err := statement.Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	var collections []*graph.Collection
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", err)
		}

		var collection graph.Collection
		var ccID sql.NullString
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
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		collection.ContractCollectionID = ccID.String
		collections = append(collections, &collection)
	}

	return collections, nil
}
