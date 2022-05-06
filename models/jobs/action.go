package jobs

type Action struct {
	ID         string `json:"id"`
	ChainID    uint64 `json:"chain_id"`
	Address    string `json:"address"`
	TokenID    string `json:"token_id"`
	ActionType string `json:"action_type"`
	Height     uint64 `json:"height"`
	Status     string `json:"status"`
	Data       []byte `json:"data"`
}
