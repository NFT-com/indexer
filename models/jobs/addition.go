package jobs

type Addition struct {
	ID      string `json:"id"`
	ChainID string `json:"chain_id"`
	Address string `json:"address"`
	TokenID string `json:"token_id"`
	Height  uint64 `json:"height"`
	Data    []byte `json:"data"`
	Status  string `json:"status"`
}
