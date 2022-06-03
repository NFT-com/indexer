package inputs

type Addition struct {
	NodeURL         string `json:"node_url"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	Standard        string `json:"standard"`
	Owner           string `json:"owner"`
	Number          uint   `json:"number"`
}
