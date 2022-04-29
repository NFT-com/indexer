package jobs

// Discovery is a job that discovers new NFTs on a blockchain.
type Discovery struct {
	ID           string   `json:"id"`
	ChainURL     string   `json:"chain_url"`
	ChainID      string   `json:"chain_id"`
	ChainType    string   `json:"chain_type"`
	BlockNumber  uint64   `json:"block_number"`
	Addresses    []string `json:"addresses"`
	StandardType string   `json:"standard_type"`
	Status       Status   `json:"status"`
}
