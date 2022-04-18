package chain

type Standard struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Events []EventType `json:"events"`
}
