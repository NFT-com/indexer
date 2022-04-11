package chain

import (
	"time"
)

type Config struct {
	ChainURL     string
	ChainType    string
	StandardType string
	Contract     string
	EventType    string
	StartIndex   string
	Batch        int64
	BatchDelay   time.Duration
}
