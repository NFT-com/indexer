package chain

type NFT struct {
	ID         string `json:"id"`
	Collection string `json:"collection"`
	TokenID    string `json:"token_id"`
	Owner      string `json:"owner"`
	Rarity     string `json:"rarity"`
}
