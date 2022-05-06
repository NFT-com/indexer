package creator

import (
	"github.com/NFT-com/indexer/config/params"
)

var DefaultConfig = Config{
	ChainID:      params.ChainEthereum,
	PendingLimit: 1000,
	HeightRange:  10,
}

type Config struct {
	ChainID      uint64
	PendingLimit uint
	HeightRange  uint
}

type Option func(*Config)

func WithChainID(chain uint64) Option {
	return func(cfg *Config) {
		cfg.ChainID = chain
	}
}

func WithPendingLimit(limit uint) Option {
	return func(cfg *Config) {
		cfg.PendingLimit = limit
	}
}

func WithHeightRange(hrange uint) Option {
	return func(cfg *Config) {
		cfg.HeightRange = hrange
	}
}
