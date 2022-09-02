package pipeline

import (
	"time"

	"github.com/NFT-com/indexer/config/params"
)

var DefaultCreationConfig = CreationConfig{
	ChainID:       params.ChainEthereum,
	CheckInterval: 2 * time.Second,
	AddressLimit:  10,
	HeightLimit:   10,
}

type CreationConfig struct {
	ChainID       uint64        // what chain ID we are watching
	CheckInterval time.Duration // how often to check for new combinations
	AddressLimit  uint          // how many addresses can be included in a single job
	HeightLimit   uint          // how many heights can be included in a single job
}

type Option func(*CreationConfig)

func WithChainID(chain uint64) Option {
	return func(cfg *CreationConfig) {
		cfg.ChainID = chain
	}
}

func WithCheckInterval(interval time.Duration) Option {
	return func(cfg *CreationConfig) {
		cfg.CheckInterval = interval
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
