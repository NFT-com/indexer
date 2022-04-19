package jobs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

// CreateDiscoveryJob creates a new discovery job and returns it.
func (c *Handler) CreateDiscoveryJob(job jobs.Discovery) (*jobs.Discovery, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err := c.store.CreateDiscoveryJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not create discovery job: %w", err)
	}

	jobList := []jobs.Discovery{job}
	err = c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, broadcaster.CreateStatusValue, jobList)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast message: %w", err)
	}

	return &job, nil
}

// CreateDiscoveryJobs creates a new discovery jobs.
func (c *Handler) CreateDiscoveryJobs(jobList []jobs.Discovery) error {
	for i := range jobList {
		jobList[i].ID = uuid.New().String()
		jobList[i].Status = jobs.StatusCreated
	}

	err := c.store.CreateDiscoveryJobs(jobList)
	if err != nil {
		return fmt.Errorf("could not create discovery jobs: %w", err)
	}

	err = c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, broadcaster.CreateStatusValue, jobList)
	if err != nil {
		return fmt.Errorf("could not broadcast message: %w", err)
	}

	return nil
}

// ListDiscoveryJobs returns a list of discovery jobs given the status. Empty string status returns all jobs.
func (c *Handler) ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error) {
	jobs, err := c.store.DiscoveryJobs(status)
	if err != nil {
		return nil, fmt.Errorf("could not get discovery jobs: %w", err)
	}

	return jobs, nil
}

func (c *Handler) GetDiscoveryJob(id string) (*jobs.Discovery, error) {
	job, err := c.store.DiscoveryJob(id)
	if err != nil {
		return nil, fmt.Errorf("could not get discovery job: %w", err)
	}

	return job, nil
}

// GetHighestBlockNumberDiscoveryJob returns the highest block number parsing block with the specified elements.
func (c *Handler) GetHighestBlockNumberDiscoveryJob(chainURL, chainType string, addresses []string, standardType, eventType string) (*jobs.Discovery, error) {
	job, err := c.store.HighestBlockNumberDiscoveryJob(chainURL, chainType, addresses, standardType, eventType)
	if err != nil {
		return nil, fmt.Errorf("could not get highest block number discovery job: %w", err)
	}

	return job, nil
}

// UpdateDiscoveryJobStatus updates the discovery job status.
func (c *Handler) UpdateDiscoveryJobStatus(id string, newStatus jobs.Status) error {
	job, err := c.store.DiscoveryJob(id)
	if err != nil {
		return fmt.Errorf("could not get discovery job: %w", err)
	}

	err = c.validateStatusSwitch(job.Status, newStatus)
	if err != nil {
		return fmt.Errorf("could not validate new job status: %w", err)
	}

	err = c.store.UpdateDiscoveryJobStatus(id, newStatus)
	if err != nil {
		return fmt.Errorf("could not update job state: %w", err)
	}
	job.Status = newStatus

	jobList := []jobs.Discovery{*job}
	err = c.BroadcastMessage(broadcaster.DiscoveryHandlerValue, broadcaster.UpdateStatusValue, jobList)
	if err != nil {
		return fmt.Errorf("could not broadcast message: %w", err)
	}

	return nil
}