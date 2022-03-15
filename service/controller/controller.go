package controller

import (
	"encoding/json"

	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/job"
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

func (c *Controller) CreateDiscoveryJob(discovery job.Discovery) (*job.Discovery, error) {
	discovery.ID = job.ID(uuid.New().String())
	discovery.Status = job.StatusCreated

	if err := c.jobsStore.CreateDiscoveryJob(discovery); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, discovery); err != nil {
		return nil, err
	}

	return &discovery, nil
}

func (c *Controller) ListDiscoveryJobs(status job.Status) ([]job.Discovery, error) {
	jobs, err := c.jobsStore.DiscoveryJobs(status)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Controller) GetDiscoveryJob(jobID job.ID) (*job.Discovery, error) {
	discovery, err := c.jobsStore.DiscoveryJob(jobID)
	if err != nil {
		return nil, err
	}

	return discovery, nil
}

func (c *Controller) UpdateDiscoveryJobState(jobID job.ID, jobStatus job.Status) error {
	discoveryJob, err := c.jobsStore.DiscoveryJob(jobID)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(discoveryJob.Status, jobStatus)
	if err != nil {
		return err
	}

	err = c.jobsStore.UpdateDiscoveryJobState(jobID, jobStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueDiscoveryJob(jobID job.ID) (*job.Discovery, error) {
	newJob, err := c.jobsStore.DiscoveryJob(jobID)
	if err != nil {
		return nil, err
	}

	newJob.ID = job.ID(uuid.New().String())
	newJob.Status = job.StatusCreated

	if err := c.jobsStore.CreateDiscoveryJob(*newJob); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, newJob); err != nil {
		return nil, err
	}

	return newJob, nil
}

func (c *Controller) CreateParsingJob(parsing job.Parsing) (*job.Parsing, error) {
	parsing.ID = job.ID(uuid.New().String())
	parsing.Status = job.StatusCreated

	if err := c.jobsStore.CreateParsingJob(parsing); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.ParsingHandlerValue, parsing); err != nil {
		return nil, err
	}

	return &parsing, nil
}

func (c *Controller) ListParsingJobs(status job.Status) ([]job.Parsing, error) {
	jobs, err := c.jobsStore.ParsingJobs(status)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Controller) GetParsingJob(jobID job.ID) (*job.Parsing, error) {
	parsing, err := c.jobsStore.ParsingJob(jobID)
	if err != nil {
		return nil, err
	}

	return parsing, nil
}

func (c *Controller) UpdateParsingJobState(jobID job.ID, jobStatus job.Status) error {
	parsingJob, err := c.jobsStore.ParsingJob(jobID)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(parsingJob.Status, jobStatus)
	if err != nil {
		return err
	}

	err = c.jobsStore.UpdateParsingJobState(jobID, jobStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueParsingJob(jobID job.ID) (*job.Parsing, error) {
	newJob, err := c.jobsStore.ParsingJob(jobID)
	if err != nil {
		return nil, err
	}

	newJob.ID = job.ID(uuid.New().String())
	newJob.Status = job.StatusCreated

	if err := c.jobsStore.CreateParsingJob(*newJob); err != nil {
		return nil, err
	}

	if err := c.BroadcastMessage(broadcaster.ParsingHandlerValue, newJob); err != nil {
		return nil, err
	}

	return newJob, nil
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

func (c *Controller) ValidateStatusSwitch(currentStatus, newStatus job.Status) error {
	switch currentStatus {
	case job.StatusCreated:
		if newStatus != job.StatusCanceled && newStatus != job.StatusQueued {
			return ErrJobStateCannotBeChanged
		}
	case job.StatusQueued:
		if newStatus != job.StatusCanceled && newStatus != job.StatusProcessing {
			return ErrJobStateCannotBeChanged
		}
	case job.StatusProcessing:
		if newStatus != job.StatusFinished && newStatus != job.StatusFailed {
			return ErrJobStateCannotBeChanged
		}
	case job.StatusCanceled, job.StatusFinished, job.StatusFailed:
		return ErrJobStateCannotBeChanged
	}

	return nil
}
