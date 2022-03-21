package request

type Discovery struct {
	ChainURL      string   `json:"chain_url" validate:"required"`
	ChainType     string   `json:"chain_type" validate:"required"`
	BlockNumber   string   `json:"block_number" validate:"required,numeric"`
	Addresses     []string `json:"addresses" validate:"required,dive,eth_addr"`
	InterfaceType string   `json:"interface_type" validate:"required"`
}

type Parsing struct {
	ChainURL      string `json:"chain_url" validate:"required"`
	ChainType     string `json:"chain_type" validate:"required"`
	BlockNumber   string `json:"block_number" validate:"required,numeric"`
	Address       string `json:"address" validate:"required,eth_addr"`
	InterfaceType string `json:"interface_type" validate:"required"`
	EventType     string `json:"event_type" validate:"required"`
}

type Status struct {
	Status string `json:"status" validate:"required"`
}
