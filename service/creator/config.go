package creator

import (
	"github.com/NFT-com/indexer/config/params"
)

var DefaultConfig = Config{
	ChainID:     params.ChainEthereum,
	BatchSize:   100,
	HeightLimit: 10,
}

type Config struct {
	ChainID     uint64 // what chain ID we are watching
	BatchSize   uint   // how many jobs to create per combination per iteration
	HeightLimit uint   // how many heights can be included in a single job
}

type Option func(*Config)

func WithChainID(chain uint64) Option {
	return func(cfg *Config) {
		cfg.ChainID = chain
	}
}

func WithBatchSize(size uint) Option {
	return func(cfg *Config) {
		cfg.BatchSize = size
	}
}

func WithHeightLimit(limit uint) Option {
	return func(cfg *Config) {
		cfg.HeightLimit = limit
	}
}
