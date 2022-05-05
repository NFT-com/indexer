package graph

type Network struct {
	ID          string `json:"id"`
	ChainID     string `json:"chain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
}
