package postgres

type chain struct {
	ID          uint64 `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Symbol      string `json:"symbol" bson:"symbol"`
	NetworkID   string `json:"network_id" bson:"network_id"`
	ChainID     string `json:"chain_id" bson:"chain_id"`
}
