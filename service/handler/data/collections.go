package data

import (
	"github.com/NFT-com/indexer/models/chain"
)

func (c *Handler) CreateCollection(collection chain.Collection) (*chain.Collection, error) {
	if err := c.store.CreateCollection(collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func (c *Handler) ListCollections() ([]chain.Collection, error) {
	collections, err := c.store.Collections()
	if err != nil {
		return nil, err
	}

	return collections, nil
}
