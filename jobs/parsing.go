package jobs

// Parsing is a job that parses an NFT's data from block data.
type Parsing struct {
	ID           string   `json:"id"`
	ChainURL     string   `json:"chain_url"`
	ChainID      string   `json:"chain_id"`
	ChainType    string   `json:"chain_type"`
	BlockNumber  string   `json:"block_number"`
	Addresses    []string `json:"addresses"`
	StandardType string   `json:"standard_type"`
	EventType    string   `json:"event_type"`
	Status       Status   `json:"status"`
}
