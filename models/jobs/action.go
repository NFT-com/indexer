package jobs

type Action struct {
	ID         string `json:"id"`
	NetworkID  string `json:"network_id"`
	Address    string `json:"address"`
	TokenID    string `json:"token_id"`
	ActionType string `json:"action_type"`
	Height     uint64 `json:"height"`
	Data       []byte `json:"data"`
	Status     string `json:"status"`
}
