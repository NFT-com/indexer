package controller

import (
	"errors"
)

var (
	ErrJobStateCannotBeChanged = errors.New("job state cannot be changed")
)
