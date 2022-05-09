package graph

type Network struct {
	ID          string `json:"id"`
	ChainID     uint64 `json:"chain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
}
