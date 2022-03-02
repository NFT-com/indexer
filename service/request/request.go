package request

type Discovery struct {
	ID            string   `json:"id"`
	ChainURL      string   `json:"chain_url"`
	ChainType     string   `json:"chain_type"`
	BlockNumber   string   `json:"block_number"`
	Addresses     []string `json:"addresses"`
	InterfaceType string   `json:"interface_type"`
	Status        string   `json:"status"`
}

type Parsing struct {
	ID            string `json:"id"`
	ChainURL      string `json:"chain_url"`
	ChainType     string `json:"chain_type"`
	BlockNumber   string `json:"block_number"`
	Address       string `json:"address"`
	InterfaceType string `json:"interface_type"`
	EventType     string `json:"event_type"`
	Status        string `json:"status"`
}
