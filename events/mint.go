package events

import (
	"time"
)

type Mint struct {
	ID              string    `json:"id"`
	CollectionID    string    `json:"collection_id"`
	Block           string    `json:"block"`
	TransactionHash string    `json:"transaction_hash"`
	Owner           string    `json:"owner"`
	EmittedAt       time.Time `json:"emitted_at"`
}
