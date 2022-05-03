package api

type CreateActionJob struct {
	ChainID    string `json:"chain_id" validate:"required"`
	Address    string `json:"address" validate:"eth_addr"`
	TokenID    string `json:"token_id" validate:"required"`
	ActionType string `json:"action_type" validate:"required"`
	Height     uint64 `json:"height" validate:"numeric"`
	Data       []byte `json:"data" validate:"required"`
}

type UpdateActionJob struct {
	Status string `json:"status" validate:"oneof:queued processing failed finished"`
}
