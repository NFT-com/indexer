package postgres

import (
	"encoding/json"
	"time"
)

type event struct {
	ID              string          `json:"id" bson:"id"`
	ChainID         uint64          `json:"chain_id" bson:"chain_id"`
	CollectionID    uint64          `json:"collection_id" bson:"collection_id"`
	Block           uint64          `json:"block" bson:"block"`
	TransactionHash string          `json:"transaction_hash" bson:"transaction_hash"`
	Type            string          `json:"type" bson:"type"`
	Data            json.RawMessage `json:"data" bson:"data"`
	EmittedAt       time.Time       `json:"emitted_at" bson:"emitted_at"`
}
