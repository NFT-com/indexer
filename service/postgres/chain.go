package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) CreateChain(chain chain.Chain) error {
	_, err := s.sqlBuilder.
		Insert(chainTableName).
		Columns(chainTableColumns...).
		Values(chain.ID, chain.ChainID, chain.Name, chain.Description, chain.Symbol).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create chain: %w", err)
	}

	return nil
}

func (s *Store) Chains() ([]chain.Chain, error) {
	result, err := s.sqlBuilder.
		Select(chainTableColumns...).
		From(chainTableName).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not list chain: %w", err)
	}

	chains := make([]chain.Chain, 0)
	for result.Next() && result.Err() == nil {
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

		chains = append(chains, chain)
	}

	return chains, nil
}

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
