package jobs

type Addition struct {
	ID              string `json:"id"`
	ChainID         uint64 `json:"chain_id"`
	CollectionID    string `json:"collection_id"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	TokenStandard   string `json:"token_standard"`
	OwnerAddress    string `json:"owner_address"`
	TokenCount      uint   `json:"token_count"`
}

func (a Addition) NFTID() string {
	return NFTID(a.ChainID, a.ContractAddress, a.TokenID)
}
