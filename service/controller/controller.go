package controller

import (
	"encoding/json"

	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

type Controller struct {
	jobsStore   JobsStore
	broadcaster *melody.Melody
}

func NewController(jobsStore JobsStore, broadcaster *melody.Melody) *Controller {
	c := Controller{
		jobsStore:   jobsStore,
		broadcaster: broadcaster,
	}

	return &c
}

func (c *Controller) CreateDiscoveryJob(job jobs.Discovery) (*jobs.Discovery, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	if err := c.jobsStore.CreateDiscoveryJob(job); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Controller) ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error) {
	jobs, err := c.jobsStore.DiscoveryJobs(status)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Controller) GetDiscoveryJob(id string) (*jobs.Discovery, error) {
	job, err := c.jobsStore.DiscoveryJob(id)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (c *Controller) UpdateDiscoveryJobState(id string, status jobs.Status) error {
	job, err := c.jobsStore.DiscoveryJob(id)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(job.Status, status)
	if err != nil {
		return err
	}

	err = c.jobsStore.UpdateDiscoveryJobState(id, status)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueDiscoveryJob(id string) (*jobs.Discovery, error) {
	job, err := c.jobsStore.DiscoveryJob(id)
	if err != nil {
		return nil, err
	}

	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	if err := c.jobsStore.CreateDiscoveryJob(*job); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (c *Controller) CreateParsingJob(job jobs.Parsing) (*jobs.Parsing, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	if err := c.jobsStore.CreateParsingJob(job); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.ParsingHandlerValue, job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Controller) ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	jobs, err := c.jobsStore.ParsingJobs(status)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Controller) GetParsingJob(id string) (*jobs.Parsing, error) {
	job, err := c.jobsStore.ParsingJob(id)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (c *Controller) UpdateParsingJobState(id string, status jobs.Status) error {
	job, err := c.jobsStore.ParsingJob(id)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(job.Status, status)
	if err != nil {
		return err
	}

	err = c.jobsStore.UpdateParsingJobState(id, status)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueParsingJob(id string) (*jobs.Parsing, error) {
	job, err := c.jobsStore.ParsingJob(id)
	if err != nil {
		return nil, err
	}

	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	if err := c.jobsStore.CreateParsingJob(*job); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.ParsingHandlerValue, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (c *Controller) BroadcastMessage(handler string, message interface{}) error {
	rawMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := c.broadcaster.BroadcastBinaryFilter(rawMessage, func(session *melody.Session) bool {
		keys := broadcaster.Keys(session.Keys)

		return keys.HasHandler(handler)
	}); err != nil {
		return err
	}

	return nil
}

func (c *Controller) ValidateStatusSwitch(currentStatus, newStatus jobs.Status) error {
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
