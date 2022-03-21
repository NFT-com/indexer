package handler

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

func (c *Handler) CreateDiscoveryJob(job jobs.Discovery) (*jobs.Discovery, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err := c.store.CreateDiscoveryJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not create discovery job: %v", err)
	}

	err = c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, job)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast message: %v", err)
	}

	return &job, nil
}

func (c *Handler) ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error) {
	jobs, err := c.store.DiscoveryJobs(status)
	if err != nil {
		return nil, fmt.Errorf("could not get discovery jobs: %v", err)
	}

	return jobs, nil
}

func (c *Handler) GetDiscoveryJob(jobID jobs.ID) (*jobs.Discovery, error) {
	job, err := c.store.DiscoveryJob(jobID)
	if err != nil {
		return nil, fmt.Errorf("could not get discovery job: %v", err)
	}

	return job, nil
}

func (c *Handler) UpdateDiscoveryJobState(jobID jobs.ID, newStatus jobs.Status) error {
	job, err := c.store.DiscoveryJob(jobID)
	if err != nil {
		return fmt.Errorf("could not get discovery job: %v", err)
	}

	err = c.ValidateStatusSwitch(job.Status, newStatus)
	if err != nil {
		return fmt.Errorf("could not validate new job status: %v", err)
	}

	err = c.store.UpdateDiscoveryJobState(jobID, newStatus)
	if err != nil {
		return fmt.Errorf("could not update job state: %v", err)
	}

	return nil
}

func (c *Handler) RequeueDiscoveryJob(jobID jobs.ID) (*jobs.Discovery, error) {
	job, err := c.store.DiscoveryJob(jobID)
	if err != nil {
		return nil, fmt.Errorf("could not get discovery job: %v", err)
	}

	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err = c.store.CreateDiscoveryJob(*job)
	if err != nil {
		return nil, fmt.Errorf("could not create discovery job: %v", err)
	}

	err = c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, job)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast message: %v", err)
	}

	return job, nil
}
