package logs

import (
	"time"
)

type Entry struct {
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
	URI                  string     `json:"uri"`
	EmittedAt            time.Time  `json:"emitted_at"`
}
