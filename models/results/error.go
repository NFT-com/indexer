package results

import (
	"errors"
)

var ErrTokenNotFound = errors.New("token not found")

type Error struct {
	Message string `json:"errorMessage"`
	Type    string `json:"errorType"`
}

func (e Error) Error() string {
	return e.Message
}
