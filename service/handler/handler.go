package handler

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

// Handler represents the business handler.
type Handler struct {
	store       JobsStore
	broadcaster *melody.Melody
}

// New returns a new business Handler.
func New(store JobsStore, broadcaster *melody.Melody) *Handler {
	c := Handler{
		store:       store,
		broadcaster: broadcaster,
	}

	return &c
}

// BroadcastMessage broadcast a message to the handlers.
func (c *Handler) BroadcastMessage(handler string, message interface{}) error {
	rawMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %w", err)
	}

	err = c.broadcaster.BroadcastBinaryFilter(rawMessage, func(session *melody.Session) bool {
		return broadcaster.HasHandler(session.Keys, handler)
	})
	if err != nil {
		return fmt.Errorf("could not broadcast message: %w", err)
	}

	return nil
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
