package handler

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/jobs"
)

// CreateAdditionJob creates a new addition job and returns it.
func (c *Handler) CreateAdditionJob(job jobs.Addition) (*jobs.Addition, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err := c.store.CreateAdditionJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not create addition job: %w", err)
	}

	return &job, nil
}

// CreateAdditionJobs creates a new addition jobs.
func (c *Handler) CreateAdditionJobs(jobList []jobs.Addition) error {
	for i := range jobList {
		jobList[i].ID = uuid.New().String()
		jobList[i].Status = jobs.StatusCreated
	}

	err := c.store.CreateAdditionJobs(jobList)
	if err != nil {
		return fmt.Errorf("could not create discovery jobs: %w", err)
	}

	return nil
}

// ListAdditionJobs returns a list of addition jobs given the status. Empty string status returns all jobs.
func (c *Handler) ListAdditionJobs(status jobs.Status) ([]jobs.Addition, error) {
	jobs, err := c.store.AdditionJobs(status)
	if err != nil {
		return nil, fmt.Errorf("could not get addition jobs: %w", err)
	}

	return jobs, nil
}

// GetAdditionJob returns an addition job given the id.
func (c *Handler) GetAdditionJob(id string) (*jobs.Addition, error) {
	job, err := c.store.AdditionJob(id)
	if err != nil {
		return nil, fmt.Errorf("could not get addition job: %w", err)
	}

	return job, nil
}

// UpdateAdditionJobStatus updates the addition job status.
func (c *Handler) UpdateAdditionJobStatus(id string, newStatus jobs.Status) error {
	job, err := c.store.AdditionJob(id)
	if err != nil {
		return fmt.Errorf("could not get addition job: %w", err)
	}

	err = c.validateStatusSwitch(job.Status, newStatus)
	if err != nil {
		return fmt.Errorf("could not validade new job status: %w", err)
	}

	err = c.store.UpdateAdditionJobStatus(id, newStatus)
	if err != nil {
		return fmt.Errorf("could not update job state: %w", err)
	}
	job.Status = newStatus

	return nil
}
