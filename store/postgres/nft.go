package postgres

import "encoding/json"

type nft struct {
	ID              uint64          `json:"id" bson:"id"`
	CollectionID    uint64          `json:"collection_id" bson:"collection_id"`
	NFTCollectionID uint            `json:"nft_collection_id" bson:"nft_collection_id"`
	Name            string          `json:"name" bson:"name"`
	URI             string          `json:"uri" bson:"uri"`
	Data            json.RawMessage `json:"data" bson:"data"`
}
