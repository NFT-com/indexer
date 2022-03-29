package handler

import (
	"errors"
)

var (
	errJobStateCannotBeChanged = errors.New("job state cannot be changed")
)
