package jobs

import (
	"github.com/NFT-com/indexer/models/id"
)

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
	return id.NFT(a.ChainID, a.ContractAddress, a.TokenID)
}

func (a Addition) TraitID(index uint) string {
	return id.Trait(a.ChainID, a.ContractAddress, a.TokenID, index)
}
