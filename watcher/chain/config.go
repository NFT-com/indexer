package chain

import (
	"math/big"
	"time"
)

type Indexes map[string]CollectionIndexes
type CollectionIndexes map[string]EventTypesIndexes
type EventTypesIndexes map[string]*big.Int

type Config struct {
	ChainURL      string
	ChainID       string
	ChainType     string
	Contracts     []string
	Standards     map[string][]string
	EventTypes    map[string][]string
	StartingBlock *big.Int
	StartIndexes  Indexes
	Batch         int64
	BatchDelay    time.Duration
}
