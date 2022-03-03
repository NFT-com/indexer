package controller

import (
	"encoding/json"

	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

// FIXME: What should I call this package?

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

func (c *Controller) GetDiscoveryJob(jobID string) (*job.Discovery, error) {
	discovery, err := c.jobsStore.DiscoveryJob(jobID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (c *Controller) UpdateDiscoveryJobState(jobID string, jobStatus job.Status) error {
	discoveryJob, err := c.jobsStore.DiscoveryJob(jobID)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(job.Status, jobStatus)
	if err != nil {
		return err
	}

	err = c.jobsStore.UpdateDiscoveryJobState(jobID, jobStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueDiscoveryJob(jobID string) (*job.Discovery, error) {
	newJob, err := c.jobsStore.DiscoveryJob(jobID)
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

func (c *Controller) GetParsingJob(jobID string) (*job.Parsing, error) {
	parsing, err := c.jobsStore.ParsingJob(jobID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (c *Controller) UpdateParsingJobState(jobID string, jobStatus job.Status) error {
	parsingJob, err := c.jobsStore.ParsingJob(jobID)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(job.Status, jobStatus)
	if err != nil {
		return err
	}

	err = c.jobsStore.UpdateParsingJobState(jobID, jobStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueParsingJob(jobID string) (*job.Parsing, error) {
	newJob, err := c.jobsStore.ParsingJob(jobID)
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
