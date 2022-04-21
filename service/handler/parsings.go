package handler

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/jobs"
)

// CreateParsingJob creates a new parsing job and returns it.
func (c *Handler) CreateParsingJob(job jobs.Parsing) (*jobs.Parsing, error) {
	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err := c.store.CreateParsingJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not create parsing job: %w", err)
	}

	return &job, nil
}

// CreateParsingJobs creates a new parsing jobs.
func (c *Handler) CreateParsingJobs(jobList []jobs.Parsing) error {
	for i := range jobList {
		jobList[i].ID = uuid.New().String()
		jobList[i].Status = jobs.StatusCreated
	}

	err := c.store.CreateParsingJobs(jobList)
	if err != nil {
		return fmt.Errorf("could not create parsing jobs: %w", err)
	}

	return nil
}

// ListParsingJobs returns a list of parsing jobs given the status. Empty string status returns all jobs.
func (c *Handler) ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	jobs, err := c.store.ParsingJobs(status)
	if err != nil {
		return nil, fmt.Errorf("could not get parsing jobs: %w", err)
	}

	return jobs, nil
}

// GetParsingJob returns a parsing job given the id.
func (c *Handler) GetParsingJob(id string) (*jobs.Parsing, error) {
	job, err := c.store.ParsingJob(id)
	if err != nil {
		return nil, fmt.Errorf("could not get parsing job: %w", err)
	}

	return job, nil
}

// GetHighestBlockNumberParsingJob returns the latest parsing block with the specified elements.
func (c *Handler) GetHighestBlockNumberParsingJob(chainURL, chainType, address, standardType, eventType string) (*jobs.Parsing, error) {
	job, err := c.store.HighestBlockNumberParsingJob(chainURL, chainType, address, standardType, eventType)
	if err != nil {
		return nil, fmt.Errorf("could not get highest block number parsing job: %w", err)
	}

	return job, nil
}

// UpdateParsingJobStatus updates the parsing job status.
func (c *Handler) UpdateParsingJobStatus(id string, newStatus jobs.Status) error {
	job, err := c.store.ParsingJob(id)
	if err != nil {
		return fmt.Errorf("could not get parsing job: %w", err)
	}

	err = c.validateStatusSwitch(job.Status, newStatus)
	if err != nil {
		return fmt.Errorf("could not validate new job status: %w", err)
	}

	err = c.store.UpdateParsingJobStatus(id, newStatus)
	if err != nil {
		return fmt.Errorf("could not update job state: %w", err)
	}
	job.Status = newStatus

	return nil
}
