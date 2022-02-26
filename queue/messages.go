package queue

type DiscoveryJob struct {
	ChainURL   string   `json:"chain_url"`
	ChainType  string   `json:"chain_type"`
	StartIndex string   `json:"start_index"`
	EndIndex   string   `json:"end_index"`
	Contracts  []string `json:"contracts"`
}

type ParseJob struct {
	ID              string                 `json:"id"`
	NetworkID       string                 `json:"network_id"`
	ChainID         string                 `json:"chain_id"`
	Block           uint64                 `json:"block"`
	TransactionHash string                 `json:"transaction_hash"`
	Address         string                 `json:"address"`
	Type            string                 `json:"type"`
	Data            map[string]interface{} `json:"data"`
}
