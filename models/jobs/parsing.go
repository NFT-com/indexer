package jobs

// Parsing is a job that parses an NFT's data from block data.
type Parsing struct {
	ID                string   `json:"id"`
	ChainID           uint64   `json:"chain_id"`
	StartHeight       uint64   `json:"start_height"`
	EndHeight         uint64   `json:"end_height"`
	ContractAddresses []string `json:"contract_addresses"`
	EventHashes       []string `json:"event_hashes"`
}
