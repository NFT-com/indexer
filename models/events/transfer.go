package events

import (
	"time"
)

type Transfer struct {
	ID                string    `json:"id"`
	ChainID           string    `json:"chain_id"`
	CollectionAddress string    `json:"collection_address"`
	TokenID           string    `json:"token_id"`
	BlockNumber       uint64    `json:"block_number"`
	TransactionHash   string    `json:"transaction_hash"`
	EventIndex        uint      `json:"event_index"`
	FromAddress       string    `json:"from_address"`
	ToAddress         string    `json:"to_address"`
	EmittedAt         time.Time `json:"emitted_at"`
}
