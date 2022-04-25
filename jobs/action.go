package jobs

// Action is a job that handles different types of actions.
type Action struct {
	ID          string `json:"id"`
	ChainURL    string `json:"chain_url"`
	ChainID     string `json:"chain_id"`
	ChainType   string `json:"chain_type"`
	BlockNumber string `json:"block_number"`
	Address     string `json:"address"`
	Standard    string `json:"standard_type"`
	Event       string `json:"event_type"`
	TokenID     string `json:"token_id"`
	Type        string `json:"type"`
	Status      Status `json:"status"`
}
