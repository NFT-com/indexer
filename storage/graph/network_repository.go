package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/graph"
)

type NetworkRepository struct {
	build squirrel.StatementBuilderType
}

func NewNetworkRepository(db *sql.DB) *NetworkRepository {

	cache := squirrel.NewStmtCache(db)
	c := NetworkRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &c
}

func (n *NetworkRepository) Retrieve(chainID string) (*graph.Network, error) {

	result, err := n.build.
		Select(ColumnsNetworks...).
		From(TableNetworks).
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

	var network graph.Network
	err = result.Scan(
		&network.ID,
		&network.ChainID,
		&network.Name,
		&network.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	return &network, nil
}
