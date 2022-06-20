package jobs

import "github.com/NFT-com/indexer/models/inputs"

type Modification struct {
	ID              string `json:"id"`
	ChainID         uint64 `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	SenderAddress   string `json:"sender_address"`
	ReceiverAddress string `json:"receiver_address"`
	TokenCount      uint   `json:"token_count"`
}

func (m Modification) NFTID() string {
	return inputs.NFTID(m.ChainID, m.ContractAddress, m.TokenID)
}
