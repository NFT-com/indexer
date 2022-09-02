package jobs

import (
	"fmt"
)

type Combination struct {
	ChainID         uint64 `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	EventHash       string `json:"event_hash"`
	StartHeight     uint64 `json:"start_height"`
}

func (c Combination) Key() string {
	return fmt.Sprintf("%d/%s/%s", c.ChainID, c.ContractAddress, c.EventHash)
}
