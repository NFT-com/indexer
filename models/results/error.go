package results

import "fmt"

type Error struct {
	Message string `json:"errorMessage"`
	Type    string `json:"errorType"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
