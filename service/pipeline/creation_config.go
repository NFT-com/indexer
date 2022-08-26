package pipeline

import (
	"github.com/NFT-com/indexer/config/params"
)

var DefaultCreationConfig = CreationConfig{
	ChainID:      params.ChainEthereum,
	BatchSize:    100,
	AddressLimit: 10,
	HeightLimit:  10,
}

type CreationConfig struct {
	ChainID      uint64 // what chain ID we are watching
	BatchSize    uint   // how many jobs to create per combination per iteration
	AddressLimit uint   // how many addresses can be included in a single job
	HeightLimit  uint   // how many heights can be included in a single job
}

type Option func(*CreationConfig)

func WithChainID(chain uint64) Option {
	return func(cfg *CreationConfig) {
		cfg.ChainID = chain
	}
}

func WithBatchSize(size uint) Option {
	return func(cfg *CreationConfig) {
		cfg.BatchSize = size
	}
}

func WithAddressLimit(limit uint) Option {
	return func(cfg *CreationConfig) {
		cfg.AddressLimit = limit
	}
}

func WithHeightLimit(limit uint) Option {
	return func(cfg *CreationConfig) {
		cfg.HeightLimit = limit
	}
}
