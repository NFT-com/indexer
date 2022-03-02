package data

import (
	"github.com/NFT-com/indexer/models/chain"
)

func (c *Handler) CreateChain(chain chain.Chain) (*chain.Chain, error) {
	if err := c.store.CreateChain(chain); err != nil {
		return nil, err
	}

	return &chain, nil
}

func (c *Handler) ListChains() ([]chain.Chain, error) {
	chains, err := c.store.Chains()
	if err != nil {
		return nil, err
	}

	return chains, nil
}
