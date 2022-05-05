package events

import (
	"time"
)

type Mint struct {
	ID              string    `json:"id"`
	CollectionID    string    `json:"collection_id"`
	TokenID         string    `json:"token_id"`
	BlockNumber     uint64    `json:"block_number"`
	TransactionHash string    `json:"transaction_hash"`
	EventIndex      uint      `json:"event_index"`
	ToAddress       string    `json:"to_address"`
	EmittedAt       time.Time `json:"emitted_at"`
}
