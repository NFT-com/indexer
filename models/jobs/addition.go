package jobs

import (
	"github.com/NFT-com/indexer/models/inputs"
)

type Addition struct {
	ID              string `json:"id"`
	ChainID         uint64 `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	TokenStandard   string `json:"token_standard"`
	OwnerAddress    string `json:"owner_address"`
	TokenCount      uint   `json:"token_count"`
}

func (a Addition) NFTID() string {
	return inputs.NFTID(a.ChainID, a.ContractAddress, a.TokenID)
}
