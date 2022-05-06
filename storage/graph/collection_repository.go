package graph

import (
	"database/sql"
	"fmt"
	"strings"

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

	query := c.build.
		Select("collections.ID, collections.network_id, collections.name, collections.description, collections.symbol, collections.slug, collections.image_url, collections.website").
		From("networks, collections").
		Where("networks.chain_id = ?", chainID).
		Where("collections.network_id = networks.id").
		Where("collections.address = ?", strings.ToLower(address))

	result, err := query.Query()
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

	return &collection, nil
}

func (c *CollectionRepository) Combinations(chainID uint64) ([]*jobs.Combination, error) {

	result, err := c.build.
		Select("collections.chain_id, collections.contract_address, events.event_hash, collections.start_height").
		From("networks, collections, collections_standards, standards, standards_events, events").
		Where("networks.chain_id = ?", chainID).
		Where("collections.network_id = networks.id").
		Where("collections.id = collections_standards.collection_id").
		Where("collection_standards.standard_id = standards.id").
		Where("standards.id = standards_events.standard_id").
		Where("standard_events.event_id = events.id").
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
