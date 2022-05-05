package events

import (
	"time"
)

type Transfer struct {
	ID              string    `json:"id"`
	CollectionID    string    `json:"collection_id"`
	Block           uint64    `json:"block"`
	EventIndex      uint      `json:"event_index"`
	TransactionHash string    `json:"transaction_hash"`
	TokenID         string    `json:"token_id"`
	Amount          uint64    `json:"amount"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	EmittedAt       time.Time `json:"emitted_at"`
}
