package events

import (
	"time"
)

type Sale struct {
	ID              string    `json:"id"`
	MarketplaceID   string    `json:"marketplace_id"`
	Block           uint64    `json:"block"`
	EventIndex      uint      `json:"event_index"`
	TransactionHash string    `json:"transaction_hash"`
	Seller          string    `json:"seller"`
	Buyer           string    `json:"buyer"`
	Price           string    `json:"price"`
	EmittedAt       time.Time `json:"emitted_at"`
}
