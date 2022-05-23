package inputs

type OwnerChange struct {
	PrevOwner string `json:"old_owner"`
	NewOwner  string `json:"new_owner"`
	Number    int64  `json:"number"`
}
