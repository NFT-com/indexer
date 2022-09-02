package jobs

import (
	"fmt"
)

type Boundary struct {
	ChainID         uint64 `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	EventHash       string `json:"event_hash"`
	NextHeight      uint64 `json:"next_height"`
}

func (b Boundary) Key() string {
	return fmt.Sprintf("%d/%s/%s", b.ChainID, b.ContractAddress, b.EventHash)
}
