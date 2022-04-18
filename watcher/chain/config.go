package chain

import (
	"time"
)

type Config struct {
	ChainURL             string
	ChainID              string
	ChainType            string
	StandardType         string
	Contracts            []string
	EventType            string
	Batch                int64
	BatchDelay           time.Duration
	ContractBlockHeights map[string]string
}
