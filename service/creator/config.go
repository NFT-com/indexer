package creator

import (
	"github.com/NFT-com/indexer/config/params"
)

var DefaultConfig = Config{
	NodeURL:      "ws://127.0.0.1:8545",
	ChainID:      params.ChainEthereum,
	PendingLimit: 1000,
	HeightRange:  10,
}

type Config struct {
	NodeURL      string
	ChainID      uint64
	PendingLimit uint
	HeightRange  uint
}

type Option func(*Config)

func WithNodeURL(url string) Option {
	return func(cfg *Config) {
		cfg.NodeURL = url
	}
}

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
