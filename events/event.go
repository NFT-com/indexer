package events

import (
	"time"
)

const (
	TypeMint     = "mint"
	TypeTransfer = "transfer"
	TypeBurn     = "burn"
	TypeSell     = "sell"
)

type Event struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	ChainID     string    `json:"chain_id"`
	NetworkID   string    `json:"network_id"`
	NftID       string    `json:"nft_id"`
	Contract    string    `json:"contract"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Price       string    `json:"price"`
	EmittedAt   time.Time `json:"emitted_at"`
}

type RawEvent struct {
	ID              string    `json:"id"`
	ChainID         string    `json:"chain_id"`
	NetworkID       string    `json:"network_id"`
	BlockNumber     string    `json:"block_number"`
	BlockHash       string    `json:"block_hash"`
	Address         string    `json:"address"`
	TransactionHash string    `json:"transaction_hash"`
	EventType       string    `json:"event_type"`
	IndexData       []string  `json:"index_data"`
	Data            []byte    `json:"data"`
	EmittedAt       time.Time `json:"emitted_at"`
}
