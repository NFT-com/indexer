package inputs

type OwnerChange struct {
	NFTID    string `json:"nft_id"`
	NewOwner string `json:"new_owner"`
}
