package jobs

// Parsing is a job that parses an NFT's data from block data.
type Parsing struct {
	ID          string `json:"id"`
	ChainURL    string `json:"chain_url"`
	ChainID     string `json:"chain_id"`
	ChainType   string `json:"chain_type"`
	BlockNumber uint64 `json:"block_number"`
	Address     string `json:"address"`
	Standard    string `json:"standard_type"`
	Event       string `json:"event_type"`
	Status      Status `json:"status"`
}

// ArchivedParsing is a successful parsing job that was archived and dropped from the Parsing table.
type ArchivedParsing struct {
	ID          string `json:"id"`
	Address     string `json:"address"`
	BlockNumber string `json:"block_number"`
	EventType   string `json:"event_type"`
}
