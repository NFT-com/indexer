package chain

type Chain struct {
	ID          string `json:"id"`
	ChainID     string `json:"chain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
}
