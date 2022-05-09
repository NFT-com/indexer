package events

import (
	"time"
)

type Sale struct {
	ID                 string    `json:"id"`
	ChainID            uint64    `json:"chain_id"`
	MarketplaceAddress string    `json:"marketplace_address"`
	TokenID            string    `json:"token_id"`
	BlockNumber        uint64    `json:"block_number"`
	TransactionHash    string    `json:"transaction_hash"`
	EventIndex         uint      `json:"event_index"`
	SellerAddress      string    `json:"seller_address"`
	BuyerAddress       string    `json:"buyer_address"`
	TradePrice         string    `json:"trade_price"`
	EmittedAt          time.Time `json:"emitted_at"`
}
