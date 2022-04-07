package chain

type NFT struct {
	ID                   string  `json:"id"`
	ChainID              string  `json:"chain_id"`
	Contract             string  `json:"contract"`
	ContractCollectionID string  `json:"contract_collection_id"`
	TokenID              string  `json:"token_id"`
	Name                 string  `json:"name"`
	Image                string  `json:"image"`
	Description          string  `json:"description"`
	Owner                string  `json:"owner"`
	Traits               []Trait `json:"traits,omitempty"`
}
