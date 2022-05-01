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

type ActionType int

const (
	Addition ActionType = iota + 1
	OwnerChange
)

func (e ActionType) String() string {
	switch e {
	case Addition:
		return "addition"
	case OwnerChange:
		return "owner_change"
	default:
		return ""
	}
}

type Log struct {
	ID                   string     `json:"id"`
	ChainID              string     `json:"chain_id"`
	Contract             string     `json:"contract"`
	Block                uint64     `json:"block"`
	Standard             string     `json:"standard"`
	Event                string     `json:"event"`
	Index                uint       `json:"index"`
	TransactionHash      string     `json:"transaction_hash"`
	Type                 EventType  `json:"type"`
	NeedsActionJob       bool       `json:"needs_action_job"`
	ActionType           ActionType `json:"action_type"`
	ContractCollectionID string     `json:"contract_collection_id"`
	NftID                string     `json:"nft_id"`
	FromAddress          string     `json:"from_address"`
	ToAddress            string     `json:"to_address"`
	Price                string     `json:"price"`
	EmittedAt            time.Time  `json:"emitted_at"`
}

type RawLog struct {
	ID              string    `json:"id"`
	ChainID         string    `json:"chain_id"`
	Index           uint      `json:"index"`
	BlockNumber     uint64    `json:"block_number"`
	BlockHash       string    `json:"block_hash"`
	Address         string    `json:"address"`
	TransactionHash string    `json:"transaction_hash"`
	EventType       string    `json:"event_type"`
	IndexData       []string  `json:"index_data"`
	Data            []byte    `json:"data"`
	EmittedAt       time.Time `json:"emitted_at"`
}
