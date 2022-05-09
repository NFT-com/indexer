package graph

type Collection struct {
	ID              string `json:"id"`
	NetworkID       string `json:"network_id"`
	ContractAddress string `json:"contract_address"`
	StartHeight     uint64 `json:"start_height"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Symbol          string `json:"symbol"`
	Slug            string `json:"slug"`
	Website         string `json:"website"`
	ImageURL        string `json:"image_url"`
}

type CollectionExtra struct {
	Collection
	ChainID     string   `json:"chain_id"`
	EventHashes []string `json:"event_hashes"`
}
