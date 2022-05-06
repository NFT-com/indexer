package jobs

// Parsing is a job that parses an NFT's data from block data.
type Parsing struct {
	ID                string   `json:"id"`
	ChainID           uint64   `json:"chain_id"`
	ContractAddresses []string `json:"contract_addresses"`
	EventHashes       []string `json:"event_hashes"`
	StartHeight       uint64   `json:"start_height"`
	EndHeight         uint64   `json:"end_height"`
	Status            string   `json:"job_status"`
	Data              []byte   `json:"data"`
}
