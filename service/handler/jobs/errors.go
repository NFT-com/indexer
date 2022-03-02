package jobs

import (
	"errors"
)

var (
	errJobStateCannotBeChanged = errors.New("job state cannot be changed")
)
