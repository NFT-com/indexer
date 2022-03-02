package request

// Discoveries represents a list of discovery jobs request.
type Discoveries struct {
	Jobs []Discovery `json:"jobs" validate:"required"`
}

// Discovery represents a request to the Discovery API.
type Discovery struct {
	ChainURL     string   `json:"chain_url" validate:"required"`
	ChainID      string   `json:"chain_id" validate:"required"`
	ChainType    string   `json:"chain_type" validate:"required"`
	BlockNumber  string   `json:"block_number" validate:"required,numeric"`
	Addresses    []string `json:"addresses" validate:"required,dive,eth_addr"`
	StandardType string   `json:"standard_type" validate:"required"`
}

// Parsings represents a list of parsing jobs request.
type Parsings struct {
	Jobs []Parsing `json:"jobs" validate:"required"`
}

// Parsing represents a request to the Parsing API.
type Parsing struct {
	ChainURL     string `json:"chain_url" validate:"required"`
	ChainID      string `json:"chain_id" validate:"required"`
	ChainType    string `json:"chain_type" validate:"required"`
	BlockNumber  string `json:"block_number" validate:"required,numeric"`
	Address      string `json:"address" validate:"required,eth_addr"`
	StandardType string `json:"standard_type" validate:"required"`
	EventType    string `json:"event_type" validate:"required"`
}

// Additions represents a list of addition jobs request.
type Additions struct {
	Jobs []Addition `json:"jobs" validate:"required"`
}

// Addition represents a request to the Addition API.
type Addition struct {
	ChainURL     string `json:"chain_url" validate:"required"`
	ChainID      string `json:"chain_id" validate:"required"`
	ChainType    string `json:"chain_type" validate:"required"`
	BlockNumber  string `json:"block_number" validate:"required,numeric"`
	Address      string `json:"address" validate:"required,eth_addr"`
	StandardType string `json:"standard_type" validate:"required"`
	TokenID      string `json:"token_id" validate:"required"`
}

// Status represents the status API change request.
type Status struct {
	Status string `json:"status" validate:"required"`
}

type Chain struct {
	ChainID     string `json:"chain_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Symbol      string `json:"symbol" validate:"required"`
}

type Marketplace struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Website     string `json:"website" validate:"required"`
}

type Collection struct {
	ChainID              string `json:"chain_id" validate:"required"`
	ContractCollectionID string `json:"contract_collection_id"`
	Address              string `json:"address" validate:"required"`
	Name                 string `json:"name" validate:"required"`
	Description          string `json:"description" validate:"required"`
	Symbol               string `json:"symbol" validate:"required"`
	Slug                 string `json:"slug" validate:"required"`
	URI                  string `json:"uri" validate:"required"`
	ImageURL             string `json:"image_url" validate:"required"`
	Website              string `json:"website" validate:"required"`
}
