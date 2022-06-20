package jobs

type Combination struct {
	ChainID         uint64 `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	EventHash       string `json:"event_hash"`
	StartHeight     uint64 `json:"last_height"`
}
