package request

type DiscoveryJob struct {
	ChainURL   string   `json:"chain_url"`
	ChainType  string   `json:"chain_type"`
	StartIndex string   `json:"start_index"`
	EndIndex   string   `json:"end_index"`
	Contracts  []string `json:"contracts"`
}
