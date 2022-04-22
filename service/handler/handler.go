package handler

import (
	"github.com/NFT-com/indexer/jobs"
)

// Handler represents the business handler.
type Handler struct {
	store JobsStore
}

// New returns a new business Handler.
func New(store JobsStore) *Handler {
	c := Handler{
		store: store,
	}

	return &c
}

// validateStatusSwitch validates if the current status is valid to change to the new status.
func (c *Handler) validateStatusSwitch(currentStatus, newStatus jobs.Status) error {
	switch currentStatus {
	case jobs.StatusCreated:
		if newStatus != jobs.StatusCanceled && newStatus != jobs.StatusQueued {
			return errJobStateCannotBeChanged
		}
	case jobs.StatusQueued:
		if newStatus != jobs.StatusCanceled && newStatus != jobs.StatusProcessing {
			return errJobStateCannotBeChanged
		}
	case jobs.StatusProcessing:
		if newStatus != jobs.StatusFinished && newStatus != jobs.StatusFailed {
			return errJobStateCannotBeChanged
		}
	case jobs.StatusCanceled, jobs.StatusFinished, jobs.StatusFailed:
		return errJobStateCannotBeChanged
	}

	return nil
}
