package events

import (
	"time"
)

type Transfer struct {
	ID                string    `json:"id"`
	ChainID           uint64    `json:"chain_id"`
	CollectionAddress string    `json:"collection_address"`
	TokenID           string    `json:"token_id"`
	BlockNumber       uint64    `json:"block_number"`
	TransactionHash   string    `json:"transaction_hash"`
	EventIndex        uint      `json:"event_index"`
	SenderAddress     string    `json:"sender_address"`
	ReceiverAddress   string    `json:"receiver_address"`
	TokenCount        uint      `json:"token_count"`
	Value             string    `json:"value"`
	EmittedAt         time.Time `json:"emitted_at"`
}
