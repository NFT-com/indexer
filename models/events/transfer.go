package events

import (
	"time"
)

type Transfer struct {
	ID              string    `json:"id"`
	CollectionID    string    `json:"collection_id"`
	Block           string    `json:"block"`
	EventIndex      uint      `json:"event_index"`
	TransactionHash string    `json:"transaction_hash"`
	TokenID         string    `json:"token_id"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	EmittedAt       time.Time `json:"emitted_at"`
}
