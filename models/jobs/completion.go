package jobs

type Completion struct {
	ID              string   `json:"id"`
	SaleID          string   `json:"sale_id"`
	ChainID         uint64   `json:"chain_id"`
	BlockNumber     uint64   `json:"block_number"`
	TransactionHash string   `json:"transaction_hash"`
	EventHashes     []string `json:"event_hashes"`
}
