package handler

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/jobs"
)

// CreateActionJob creates a new action job and returns it.
func (c *Handler) CreateActionJob(job *jobs.Action) (*jobs.Action, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err := c.store.CreateActionJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not create action job: %w", err)
	}

	return job, nil
}

// CreateActionJobs creates a new action jobs.
func (c *Handler) CreateActionJobs(jobList []*jobs.Action) error {
	for i := range jobList {
		jobList[i].ID = uuid.New().String()
		jobList[i].Status = jobs.StatusCreated
	}

	err := c.store.CreateActionJobs(jobList)
	if err != nil {
		return fmt.Errorf("could not create action jobs: %w", err)
	}

	return nil
}

// ListActionJobs returns a list of action jobs given the status. Empty string status returns all jobs.
func (c *Handler) ListActionJobs(status jobs.Status) ([]*jobs.Action, error) {
	jobs, err := c.store.ActionJobs(status)
	if err != nil {
		return nil, fmt.Errorf("could not get action jobs: %w", err)
	}

	return jobs, nil
}

// GetActionJob returns an action job given the id.
func (c *Handler) GetActionJob(id string) (*jobs.Action, error) {
	job, err := c.store.ActionJob(id)
	if err != nil {
		return nil, fmt.Errorf("could not get action job: %w", err)
	}

	return job, nil
}

// UpdateActionJobStatus updates the action job status.
func (c *Handler) UpdateActionJobStatus(id string, newStatus jobs.Status) error {
	job, err := c.store.ActionJob(id)
	if err != nil {
		return fmt.Errorf("could not get action job: %w", err)
	}

	err = c.validateStatusSwitch(job.Status, newStatus)
	if err != nil {
		return fmt.Errorf("could not validade new job status: %w", err)
	}

	err = c.store.UpdateActionJobStatus(id, newStatus)
	if err != nil {
		return fmt.Errorf("could not update job state: %w", err)
	}
	job.Status = newStatus

	return nil
}
