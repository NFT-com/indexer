package api

type CreateActionJob struct {
	ChainURL    string `json:"chain_url" validate:"required"`
	ChainID     string `json:"chain_id" validate:"required"`
	ChainType   string `json:"chain_type" validate:"required"`
	BlockNumber uint64 `json:"block_number" validate:"required,numeric"`
	Address     string `json:"address" validate:"required,eth_addr"`
	Standard    string `json:"standard_type" validate:"required"`
	Event       string `json:"event_type" validate:"required"`
	TokenID     string `json:"token_id" validate:"required"`
	ToAddress   string `json:"to_address"`
	Type        string `json:"type" validate:"required"`
}

type UpdateActionJob struct {
	Status string `json:"status" validate:"oneof:queued processing canceled failed finished"`
}
