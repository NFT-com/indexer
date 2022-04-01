package log

import (
	"time"
)

type EventType int

const (
	Mint EventType = iota + 1
	Transfer
	Burn
	Sale
)

func (e EventType) String() string {
	switch e {
	case Mint:
		return "mint"
	case Transfer:
		return "transfer"
	case Burn:
		return "burn"
	case Sale:
		return "sale"
	default:
		return ""
	}
}

type Log struct {
	ID                   string    `json:"id"`
	ChainID              string    `json:"chain_id"`
	Contract             string    `json:"contract"`
	Block                string    `json:"block"`
	Index                uint      `json:"index"`
	TransactionHash      string    `json:"transaction_hash"`
	Type                 EventType `json:"type"`
	ContractCollectionID string    `json:"contract_collection_id"`
	NftID                string    `json:"nft_id"`
	FromAddress          string    `json:"from_address"`
	ToAddress            string    `json:"to_address"`
	Price                string    `json:"price"`
	EmittedAt            time.Time `json:"emitted_at"`
}

type RawLog struct {
	ID              string    `json:"id"`
	ChainID         string    `json:"chain_id"`
	Index           uint      `json:"index"`
	BlockNumber     string    `json:"block_number"`
	BlockHash       string    `json:"block_hash"`
	Address         string    `json:"address"`
	TransactionHash string    `json:"transaction_hash"`
	EventType       string    `json:"event_type"`
	IndexData       []string  `json:"index_data"`
	Data            []byte    `json:"data"`
	EmittedAt       time.Time `json:"emitted_at"`
}
