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
	ImageURL        string `json:"image_url"`
	Website         string `json:"website"`
}

type CollectionExtra struct {
	Collection
	ChainID     string   `json:"chain_id"`
	EventHashes []string `json:"event_hashes"`
}
