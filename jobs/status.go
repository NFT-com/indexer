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

	ErrStatusNotFound = errors.New("status not found")
)

type Status string

func ParseStatus(rawStatus string) (Status, error) {
	if rawStatus == "" {
		return Status(rawStatus), nil
	}

	rawStatus = strings.ToLower(rawStatus)

	_, exists := statusMap[rawStatus]
	if !exists {
		return "", ErrStatusNotFound
	}

	return Status(rawStatus), nil
}
