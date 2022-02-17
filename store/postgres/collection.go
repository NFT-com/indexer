package postgres

type collection struct {
	ID          uint64 `json:"id" bson:"id"`
	ChainID     uint64 `json:"chain_id" bson:"chain_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Symbol      string `json:"symbol" bson:"symbol"`
	Address     string `json:"address" bson:"address"`
	ABI         string `json:"abi" bson:"abi"`
	Standard    string `json:"standard" bson:"standard"`
}
