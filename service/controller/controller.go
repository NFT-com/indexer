package controller

import (
	"encoding/json"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/job"
)

// FIXME: What should I call this package?

type Controller struct {
	discoveryJobsStore DiscoveryJobsStore
	parsingJobsStore   ParsingJobsStore
	broadcaster        *melody.Melody
}

func NewController(discoveryJobsStore DiscoveryJobsStore, parsingJobsStore ParsingJobsStore, broadcaster *melody.Melody) *Controller {
	c := Controller{
		discoveryJobsStore: discoveryJobsStore,
		parsingJobsStore:   parsingJobsStore,
		broadcaster:        broadcaster,
	}

	return &c
}

func (c *Controller) CreateDiscoveryJob(discovery job.Discovery) (job.Discovery, error) {
	discovery.ID = uuid.New().String()
	discovery.Status = job.StatusCreated

	if err := c.discoveryJobsStore.CreateDiscoveryJob(discovery); err != nil {
		return job.Discovery{}, err
	}

	rawMessage, err := json.Marshal(discovery)
	if err != nil {
		return job.Discovery{}, err
	}

	if err := c.broadcaster.BroadcastBinaryFilter(rawMessage, func(session *melody.Session) bool {
		keys := broadcaster.NewKeys(session.Keys)

		return keys.HasHandler(broadcaster.DiscoveryHandlerValue)
	}); err != nil {
		return job.Discovery{}, err
	}

	return discovery, nil
}

func (c *Controller) ListDiscoveryJobs(status job.Status) ([]job.Discovery, error) {
	jobs, err := c.discoveryJobsStore.ListDiscoveryJobs(status)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Controller) GetDiscoveryJob(jobID job.ID) (job.Discovery, error) {
	discovery, err := c.discoveryJobsStore.GetDiscoveryJob(jobID)
	if err != nil {
		return job.Discovery{}, err
	}

	return discovery, nil
}

func (c *Controller) UpdateDiscoveryJobState(jobID job.ID, jobStatus job.Status) error {
	discoveryJob, err := c.discoveryJobsStore.GetDiscoveryJob(jobID)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(discoveryJob.Status, jobStatus)
	if err != nil {
		return err
	}

	err = c.discoveryJobsStore.UpdateDiscoveryJobState(jobID, jobStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueDiscoveryJob(jobID job.ID) (job.Discovery, error) {
	newJob, err := c.discoveryJobsStore.GetDiscoveryJob(jobID)
	if err != nil {
		return job.Discovery{}, err
	}

	newJob.ID = uuid.New().String()
	newJob.Status = job.StatusCreated

	if err := c.discoveryJobsStore.CreateDiscoveryJob(newJob); err != nil {
		return job.Discovery{}, err
	}

	return newJob, nil
}

func (c *Controller) CreateParsingJob(parsing job.Parsing) (job.Parsing, error) {
	parsing.ID = uuid.New().String()
	parsing.Status = job.StatusCreated

	if err := c.parsingJobsStore.CreateParsingJob(parsing); err != nil {
		return job.Parsing{}, err
	}

	rawMessage, err := json.Marshal(parsing)
	if err != nil {
		return job.Parsing{}, err
	}

	if err := c.broadcaster.BroadcastBinaryFilter(rawMessage, func(session *melody.Session) bool {
		keys := broadcaster.NewKeys(session.Keys)

		return keys.HasHandler(broadcaster.ParsingHandlerValue)
	}); err != nil {
		return job.Parsing{}, err
	}

	return parsing, nil
}

func (c *Controller) ListParsingJobs(status job.Status) ([]job.Parsing, error) {
	jobs, err := c.parsingJobsStore.ListParsingJobs(status)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Controller) GetParsingJob(jobID job.ID) (job.Parsing, error) {
	parsing, err := c.parsingJobsStore.GetParsingJob(jobID)
	if err != nil {
		return job.Parsing{}, err
	}

	return parsing, nil
}

func (c *Controller) UpdateParsingJobState(jobID job.ID, jobStatus job.Status) error {
	parsingJob, err := c.parsingJobsStore.GetParsingJob(jobID)
	if err != nil {
		return err
	}

	err = c.ValidateStatusSwitch(parsingJob.Status, jobStatus)
	if err != nil {
		return err
	}

	err = c.parsingJobsStore.UpdateParsingJobState(jobID, jobStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) RequeueParsingJob(jobID job.ID) (job.Parsing, error) {
	newJob, err := c.parsingJobsStore.GetParsingJob(jobID)
	if err != nil {
		return job.Parsing{}, err
	}

	newJob.ID = uuid.New().String()
	newJob.Status = job.StatusCreated

	if err := c.parsingJobsStore.CreateParsingJob(newJob); err != nil {
		return job.Parsing{}, err
	}

	return newJob, nil
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
