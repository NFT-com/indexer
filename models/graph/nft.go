package graph

type NFT struct {
	ID           string `json:"id"`
	CollectionID string `json:"collection_id"`
	TokenID      string `json:"token_id"`
	Name         string `json:"name"`
	URI          string `json:"uri"`
	Image        string `json:"image"`
	Description  string `json:"description"`
	Owner        string `json:"owner"`
}
