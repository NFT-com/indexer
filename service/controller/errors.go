package controller

import (
	"errors"
)

var (
	ErrJobStateCannotBeChanged = errors.New("jobs state cannot be changed")
)
