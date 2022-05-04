package graph

// TODO: review all of the model relationships properly

type Collection struct {
	ID              string `json:"id"`
	NetworkID       string `json:"network_id"`
	ContractAddress string `json:"contract_address"`
	BaseTokenID     string `json:"base_token_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Symbol          string `json:"symbol"`
	Slug            string `json:"slug"`
	ImageURL        string `json:"image_url"`
	Website         string `json:"website"`
}
