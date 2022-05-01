package jobs

type Input struct {
	IDs        []string          `json:"ids"`
	ChainURL   string            `json:"chain_url"`
	ChainID    string            `json:"chain_id"`
	ChainType  string            `json:"chain_type"`
	StartBlock uint64            `json:"starting_block"`
	EndBlock   uint64            `json:"end_block"`
	Addresses  []string          `json:"addresses"`
	Standards  map[string]string `json:"standards"`
	EventTypes []string          `json:"event_types"`
}
