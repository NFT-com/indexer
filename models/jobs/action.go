package jobs

type Action struct {
	ID              string `json:"id"`
	ChainID         uint64 `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	ActionType      string `json:"action_type"`
	BlockHeight     uint64 `json:"block_height"`
	JobStatus       string `json:"job_status"`
	Data            []byte `json:"data"`
}
