package events

import (
	"time"
)

type Transfer struct {
	ID                string    `gorm:"column:id" json:"id"`
	ChainID           uint64    `gorm:"column:chain_id" json:"chain_id"`
	CollectionAddress string    `gorm:"column:collection_address" json:"collection_address"`
	TokenID           string    `gorm:"column:token_id" json:"token_id"`
	BlockNumber       uint64    `gorm:"column:block_number" json:"block_number"`
	TransactionHash   string    `gorm:"column:transaction_hash" json:"transaction_hash"`
	EventIndex        uint      `gorm:"column:event_index" json:"event_index"`
	SenderAddress     string    `gorm:"column:sender_address" json:"sender_address"`
	ReceiverAddress   string    `gorm:"column:receiver_address" json:"receiver_address"`
	TokenCount        uint      `gorm:"column:token_count" json:"token_count"`
	EmittedAt         time.Time `gorm:"column:emitted_at" json:"emitted_at"`
}
