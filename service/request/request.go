package request

type Discoveries struct {
	Jobs []Discovery `json:"jobs" validate:"required"`
}

// Discovery represents a request to the Discovery API.
type Discovery struct {
	ChainURL     string   `json:"chain_url" validate:"required"`
	ChainType    string   `json:"chain_type" validate:"required"`
	BlockNumber  string   `json:"block_number" validate:"required,numeric"`
	Addresses    []string `json:"addresses" validate:"required,dive,eth_addr"`
	StandardType string   `json:"standard_type" validate:"required"`
}

type Parsings struct {
	Jobs []Parsing `json:"jobs" validate:"required"`
}

// Parsing represents a request to the Parsing API.
type Parsing struct {
	ChainURL     string `json:"chain_url" validate:"required"`
	ChainType    string `json:"chain_type" validate:"required"`
	BlockNumber  string `json:"block_number" validate:"required,numeric"`
	Address      string `json:"address" validate:"required,eth_addr"`
	StandardType string `json:"standard_type" validate:"required"`
	EventType    string `json:"event_type" validate:"required"`
}

// Status represents the status API change request.
type Status struct {
	Status string `json:"status" validate:"required"`
}
