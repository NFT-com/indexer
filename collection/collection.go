package collection

type Collection struct {
	ID                   string `json:"id"`
	ChainID              string `json:"chain_id"`
	ContractCollectionID string `json:"contract_collection_id"`
	Address              string `json:"address"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Symbol               string `json:"symbol"`
	Slug                 string `json:"slug"`
	Standard             string `json:"standard"`
	URI                  string `json:"uri"`
	ImageURL             string `json:"image_url"`
}