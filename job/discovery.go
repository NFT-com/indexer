package job

type Discovery struct {
	ID            string   `json:"id"`
	ChainURL      string   `json:"chain_url"`
	ChainType     string   `json:"chain_type"`
	BlockNumber   string   `json:"block_number"`
	Addresses     []string `json:"addresses"`
	InterfaceType string   `json:"interface_type"`
	Status        Status   `json:"status"`
}
