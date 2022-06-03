package inputs

type OwnerChange struct {
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	PrevOwner       string `json:"old_owner"`
	NewOwner        string `json:"new_owner"`
	Number          uint   `json:"number"`
}
