package chain

type Collection struct {
	ID                   string `json:"id"`
	ChainID              string `json:"chain_id"`
	ContractCollectionID string `json:"contract_collection_id"`
	Address              string `json:"address"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Symbol               string `json:"symbol"`
	Slug                 string `json:"slug"`
	ImageURL             string `json:"image_url"`
	Website              string `json:"website"`
}
