package graph

type Marketplace struct {
	ID              string `json:"id"`
	NetworkID       string `json:"network_id"`
	ContractAddress string `json:"contract_address"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Website         string `json:"website"`
}
