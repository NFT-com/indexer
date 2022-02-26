package request

type ParseJob struct {
	ID              string   `json:"id"`
	NetworkID       string   `json:"network_id"`
	ChainID         string   `json:"chain_id"`
	Block           uint64   `json:"block"`
	TransactionHash string   `json:"transaction_hash"`
	Address         string   `json:"address"`
	AddressType     string   `json:"address_type"`
	Topic           string   `json:"topic"`
	IndexedData     []string `json:"indexed_data"`
	Data            []byte   `json:"data"`
}
