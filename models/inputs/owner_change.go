package inputs

type OwnerChange struct {
	PrevOwner string `json:"old_owner"`
	NewOwner  string `json:"new_owner"`
	Number    uint64 `json:"number"`
}
