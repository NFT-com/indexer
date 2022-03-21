package handler

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
)

func (c *Handler) CreateParsingJob(parsing jobs.Parsing) (*jobs.Parsing, error) {
	parsing.ID = uuid.New().String()
	parsing.Status = jobs.StatusCreated

	err := c.store.CreateParsingJob(parsing)
	if err != nil {
		return nil, fmt.Errorf("could not create parsing job: %v", err)
	}

	err = c.BroadcastMessage(broadcaster.ParsingHandlerValue, parsing)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast message: %v", err)
	}

	return &parsing, nil
}

func (c *Handler) ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	jobs, err := c.store.ParsingJobs(status)
	if err != nil {
		return nil, fmt.Errorf("could not get parsing jobs: %v", err)
	}

	return jobs, nil
}

func (c *Handler) GetParsingJob(jobID jobs.ID) (*jobs.Parsing, error) {
	job, err := c.store.ParsingJob(jobID)
	if err != nil {
		return nil, fmt.Errorf("could not get parsing job: %v", err)
	}

	return job, nil
}

func (c *Handler) UpdateParsingJobState(jobID jobs.ID, newStatus jobs.Status) error {
	job, err := c.store.ParsingJob(jobID)
	if err != nil {
		return fmt.Errorf("could not get parsing job: %v", err)
	}

	err = c.ValidateStatusSwitch(job.Status, newStatus)
	if err != nil {
		return fmt.Errorf("could not create parsing job: %v", err)
	}

	err = c.store.UpdateParsingJobState(jobID, newStatus)
	if err != nil {
		return fmt.Errorf("could not update job state: %v", err)
	}

	return nil
}

func (c *Handler) RequeueParsingJob(jobID jobs.ID) (*jobs.Parsing, error) {
	job, err := c.store.ParsingJob(jobID)
	if err != nil {
		return nil, fmt.Errorf("could not get parsing job: %v", err)
	}

	job.ID = uuid.New().String()
	job.Status = jobs.StatusCreated

	err = c.store.CreateParsingJob(*job)
	if err != nil {
		return nil, fmt.Errorf("could not create parsing job: %v", err)
	}

	err = c.BroadcastMessage(broadcaster.ParsingHandlerValue, job)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast message: %v", err)
	}

	return job, nil
}
