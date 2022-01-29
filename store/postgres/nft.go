package postgres

import (
	"encoding/json"

	"github.com/NFT-com/indexer/nft"
)

type NFT struct {
	ID       string          `json:"id" bson:"id"`
	Network  string          `json:"network" bson:"network"`
	Chain    string          `json:"chain" bson:"chain"`
	Contract string          `json:"contract" bson:"contract"`
	Owner    string          `json:"owner" bson:"owner"`
	Data     json.RawMessage `json:"data" bson:"data"`
}

func FromExternalNFT(nft *nft.NFT) (*NFT, error) {
	internalNFT := NFT{
		ID:       nft.ID,
		Network:  nft.Network,
		Chain:    nft.Chain,
		Contract: nft.Contract,
		Owner:    nft.Owner,
	}

	data, err := json.Marshal(nft.Data)
	if err != nil {
		return nil, err
	}

	internalNFT.Data = data
	return &internalNFT, nil
}
