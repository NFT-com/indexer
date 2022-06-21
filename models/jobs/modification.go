package jobs

type Modification struct {
	ID              string `json:"id"`
	ChainID         uint64 `json:"chain_id"`
	CollectionID    string `json:"collection_id"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	SenderAddress   string `json:"sender_address"`
	ReceiverAddress string `json:"receiver_address"`
	TokenCount      uint   `json:"token_count"`
}

func (m Modification) NFTID() string {
	return NFTID(m.ChainID, m.ContractAddress, m.TokenID)
}
