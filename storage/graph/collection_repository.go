package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
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

func (c *CollectionRepository) One(chainID uint64, address string) (*graph.Collection, error) {

	result, err := c.build.
		Select("collections.ID, collections.contract_address, collections.network_id, collections.name, collections.description, collections.symbol, collections.slug, collections.image_url, collections.website").
		From("networks, collections").
		Where("networks.chain_id = ?", chainID).
		Where("collections.network_id = networks.id").
		Where("LOWER(collections.contract_address) = LOWER(?)", address).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collection: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not get row: %w", result.Err())
	}
	if !result.Next() {
		return nil, sql.ErrNoRows
	}

	var collection graph.Collection
	err = result.Scan(
		&collection.ID,
		&collection.NetworkID,
		&collection.ContractAddress,
		&collection.Name,
		&collection.Description,
		&collection.Symbol,
		&collection.Slug,
		&collection.Website,
		&collection.ImageURL,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	return &collection, nil
}

func (c *CollectionRepository) Combinations(chainID uint64) ([]*jobs.Combination, error) {

	result, err := c.build.
		Select("networks.chain_id, collections.contract_address, events.event_hash, collections.start_height").
		From("networks, collections, collections_standards, standards, standards_events, events").
		Where("networks.chain_id = ?", chainID).
		Where("collections.network_id = networks.id").
		Where("collections.id = collections_standards.collection_id").
		Where("collections_standards.standard_id = standards.id").
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
