package postgres

import "encoding/json"

type event struct {
	ID              string          `json:"id" bson:"id"`
	NetworkID       uint64          `json:"network_id" bson:"network_id"`
	CollectionID    uint64          `json:"collection_id" bson:"collection_id"`
	Block           uint64          `json:"block" bson:"block"`
	TransactionHash string          `json:"transaction_hash" bson:"transaction_hash"`
	Type            string          `json:"type" bson:"type"`
	Data            json.RawMessage `json:"data" bson:"data"`
}
