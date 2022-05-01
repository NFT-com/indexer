package logs

import (
	"time"
)

type Raw struct {
	ID              string    `json:"id"`
	ChainID         string    `json:"chain_id"`
	Index           uint      `json:"index"`
	BlockNumber     uint64    `json:"block_number"`
	BlockHash       string    `json:"block_hash"`
	Address         string    `json:"address"`
	TransactionHash string    `json:"transaction_hash"`
	EventType       string    `json:"event_type"`
	IndexData       []string  `json:"index_data"`
	Data            []byte    `json:"data"`
	EmittedAt       time.Time `json:"emitted_at"`
}
