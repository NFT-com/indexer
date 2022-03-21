package handler

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

type Handler struct {
	store       JobsStore
	broadcaster *melody.Melody
}

func New(store JobsStore, broadcaster *melody.Melody) *Handler {
	c := Handler{
		store:       store,
		broadcaster: broadcaster,
	}

	return &c
}

func (c *Handler) BroadcastMessage(handler string, message interface{}) error {
	rawMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}

	err = c.broadcaster.BroadcastBinaryFilter(rawMessage, func(session *melody.Session) bool {
		keys := broadcaster.Keys(session.Keys)

		return keys.HasHandler(handler)
	})
	if err != nil {
		return fmt.Errorf("could not broadcast message: %v", err)
	}

	return nil
}

func (c *Handler) ValidateStatusSwitch(currentStatus, newStatus jobs.Status) error {
	switch currentStatus {
	case jobs.StatusCreated:
		if newStatus != jobs.StatusCanceled && newStatus != jobs.StatusQueued {
			return ErrJobStateCannotBeChanged
		}
	case jobs.StatusQueued:
		if newStatus != jobs.StatusCanceled && newStatus != jobs.StatusProcessing {
			return ErrJobStateCannotBeChanged
		}
	case jobs.StatusProcessing:
		if newStatus != jobs.StatusFinished && newStatus != jobs.StatusFailed {
			return ErrJobStateCannotBeChanged
		}
	case jobs.StatusCanceled, jobs.StatusFinished, jobs.StatusFailed:
		return ErrJobStateCannotBeChanged
	}

	return nil
}
