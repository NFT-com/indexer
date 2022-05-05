package chain

type NFT struct {
	ID          string  `json:"id"`
	ChainID     string  `json:"chain_id"`
	Contract    string  `json:"contract"`
	TokenID     string  `json:"token_id"`
	Name        string  `json:"name"`
	URI         string  `json:"uri"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Owner       string  `json:"owner"`
	Traits      []Trait `json:"traits,omitempty"`
}
