package parse

type NFT struct {
	ID      int64  `json:"nft_id"`
	Address string `json:"address"`
	Chain   string `json:"chain"`
	Network string `json:"network"`
}
