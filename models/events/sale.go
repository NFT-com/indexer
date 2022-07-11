package events

import (
	"time"

	"golang.org/x/crypto/sha3"
)

type Sale struct {
	ID                 string    `json:"id"`
	ChainID            uint64    `json:"chain_id"`
	MarketplaceAddress string    `json:"marketplace_address"`
	CollectionAddress  string    `json:"collection_address"`
	TokenID            string    `json:"token_id"`
	BlockNumber        uint64    `json:"block_number"`
	TransactionHash    string    `json:"transaction_hash"`
	EventIndex         uint      `json:"event_index"`
	SellerAddress      string    `json:"seller_address"`
	BuyerAddress       string    `json:"buyer_address"`
	TradePrice         string    `json:"trade_price"`
	PaymentToken       string    `json:"payment_token"`
	EmittedAt          time.Time `json:"emitted_at"`
}

func (s Sale) Hash() string {
	hash := sha3.Sum256([]byte(s.TransactionHash + s.SellerAddress + s.BuyerAddress))
	return string(hash[:])
}
