package job

import (
	"errors"
	"strings"
)

const (
	StatusCreated  = "created"
	StatusQueued   = "queued"
	StatusCanceled = "canceled"
	StatusFailed   = "failed"
	StatusFinished = "finished"
)

var (
	statusMap = map[string]struct{}{
		StatusCreated:  struct{}{},
		StatusQueued:   struct{}{},
		StatusCanceled: struct{}{},
		StatusFailed:   struct{}{},
		StatusFinished: struct{}{},
	}

	ErrStatusNotFound = errors.New("status not found")
)

type Status string

func ParseStatus(rawStatus string) (Status, error) {
	if rawStatus == "" {
		return Status(rawStatus), nil
	}

	_, hasStatus := statusMap[strings.ToLower(rawStatus)]
	if !hasStatus {
		return "", ErrStatusNotFound
	}

	return Status(rawStatus), nil
}
