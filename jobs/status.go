package jobs

import (
	"errors"
	"strings"
)

const (
	StatusCreated    = "created"
	StatusQueued     = "queued"
	StatusProcessing = "processing"
	StatusCanceled   = "canceled"
	StatusFailed     = "failed"
	StatusFinished   = "finished"
)

var (
	statusMap = map[string]struct{}{
		StatusCreated:    {},
		StatusQueued:     {},
		StatusProcessing: {},
		StatusCanceled:   {},
		StatusFailed:     {},
		StatusFinished:   {},
	}

	errStatusNotFound = errors.New("status not found")
)

// Status represents the job status.
type Status string

// ParseStatus returns the status type from a raw status string.
// Returns empty status if the status is empty.
func ParseStatus(rawStatus string) (Status, error) {
	if rawStatus == "" {
		return Status(rawStatus), nil
	}

	rawStatus = strings.ToLower(rawStatus)

	_, exists := statusMap[rawStatus]
	if !exists {
		return "", errStatusNotFound
	}

	return Status(rawStatus), nil
}
