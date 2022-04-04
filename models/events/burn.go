package events

import (
	"time"
)

type Burn struct {
	ID              string    `json:"id"`
	CollectionID    string    `json:"collection_id"`
	Block           string    `json:"block"`
	EventIndex      uint      `json:"event_index"`
	TransactionHash string    `json:"transaction_hash"`
	TokenID         string    `json:"token_id"`
	EmittedAt       time.Time `json:"emitted_at"`
}
