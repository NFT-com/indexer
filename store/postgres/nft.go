package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
)

type nft struct {
	ID           uuid.UUID       `json:"id" bson:"id"`
	CollectionID uint64          `json:"collection_id" bson:"collection_id"`
	TokenID      string          `json:"token_id" bson:"token_id"`
	Owner        string          `json:"owner" bson:"owner"`
	Name         string          `json:"name" bson:"name"`
	URI          string          `json:"uri" bson:"uri"`
	Rarity       uint64          `json:"rarity" bson:"rarity"`
	Data         json.RawMessage `json:"data" bson:"data"`
}
