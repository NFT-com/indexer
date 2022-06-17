package jobs

import (
	"time"
)

type Boundary struct {
	ChainID         uint64    `json:"chain_id"`
	ContractAddress string    `json:"contract_address"`
	EventHash       []string  `json:"event_hash"`
	LastHeight      uint64    `json:"last_height"`
	LastID          string    `json:"last_id"`
	UpdatedAt       time.Time `json:"updated_at"`
}
