package jobs

type Action struct {
	ID              string `json:"id"`
	ChainID         uint64 `json:"chain_id"`
	BlockHeight     uint64 `json:"block_height"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	TokenStandard   string `json:"token_standard"`
	TokenOwner      string `json:"token_owner"`
	TokenCount      uint   `json:"token_count"`
}
