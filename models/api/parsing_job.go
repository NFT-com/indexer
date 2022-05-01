package api

type CreateParsingJob struct {
	ChainURL    string `json:"chain_url" validate:"required"`
	ChainID     string `json:"chain_id" validate:"required"`
	ChainType   string `json:"chain_type" validate:"required"`
	BlockNumber uint64 `json:"block_number" validate:"required,numeric"`
	Address     string `json:"address" validate:"required,eth_addr"`
	Standard    string `json:"standard_type" validate:"required"`
	Event       string `json:"event_type" validate:"required"`
}

type UpdateParsingJob struct {
	Status string `json:"status" validate:"oneof:queued processing failed finished"`
}
