package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/chain"
)

func (s *Store) Chain(chainID string) (*chain.Chain, error) {
	result, err := s.sqlBuilder.
		Select(chainTableColumns...).
		From(chainTableName).
		Where("chain_id = ?", chainID).
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve chain: %w", errResourceNotFound)
	}

	var chain chain.Chain
	err = result.Scan(
		&chain.ID,
		&chain.ChainID,
		&chain.Name,
		&chain.Description,
		&chain.Symbol,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan chain: %w", err)
	}

	return &chain, nil
}
