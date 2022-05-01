package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type ChainRepository struct {
	build squirrel.StatementBuilderType
}

func NewChainRepository(db *sql.DB) *ChainRepository {

	cache := squirrel.NewStmtCache(db)
	c := ChainRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &c
}

func (c *ChainRepository) Retrieve(chainID string) (*graph.Chain, error) {

	result, err := c.build.
		Select(ColumnsChains...).
		From(TableChains).
		Where("chain_id = ?", chainID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not get row: %w", err)
	}
	if !result.Next() {
		return nil, sql.ErrNoRows
	}

	var chain graph.Chain
	err = result.Scan(
		&chain.ID,
		&chain.ChainID,
		&chain.Name,
		&chain.Description,
		&chain.Symbol,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	return &chain, nil
}
