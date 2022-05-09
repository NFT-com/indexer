package api

type CreateParsingJob struct {
	ChainID     string   `json:"chain_id" validate:"required"`
	Addresses   []string `json:"addresses" validate:"required"`
	EventTypes  []string `json:"event_types" validate:"required"`
	StartHeight uint64   `json:"start_height" validate:"numeric"`
	EndHeight   uint64   `json:"end_height" validate:"numeric"`
	Data        []byte   `json:"data" validate:"required"`
}

type UpdateParsingJob struct {
	Status string `json:"status" validate:"oneof:queued processing failed finished"`
}
