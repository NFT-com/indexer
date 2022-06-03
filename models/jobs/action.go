package jobs

type Action struct {
	ID          string `json:"id"`
	ChainID     uint64 `json:"chain_id"`
	ActionType  string `json:"action_type"`
	BlockHeight uint64 `json:"block_height"`
	JobStatus   string `json:"job_status"`
	InputData   []byte `json:"input_data"`
}
