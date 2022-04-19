package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

// CreateActionJob creates a new Action job.
func (s *Store) CreateActionJob(job *jobs.Action) error {

	_, err := s.build.
		Insert(actionJobsTableName).
		Columns(actionJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, job.Address, job.Standard, job.TokenID, job.Status).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create action job: %w", err)
	}

	return nil
}

// CreateActionJobs creates a batch of Action jobs.
func (s *Store) CreateActionJobs(jobs []*jobs.Action) error {

	query := s.build.
		Insert(actionJobsTableName).
		Columns(actionJobsTableColumns...)

	for _, job := range jobs {
		query = query.Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, job.Address, job.Standard, job.TokenID, job.Status)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not create action jobs: %w", err)
	}

	return nil
}

// ActionJobs returns a list of Action jobs filtered by status. Empty string status returns every job.
func (s *Store) ActionJobs(status jobs.Status) ([]*jobs.Action, error) {

	query := s.build.
		Select(actionJobsTableColumns...).
		From(actionJobsTableName)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve action job list: %w", err)
	}
	defer result.Close()

	jobList := make([]*jobs.Action, 0)
	for result.Next() && result.Err() == nil {
		var job jobs.Action
		err = result.Scan(
			&job.ID,
			&job.ChainURL,
			&job.ChainID,
			&job.ChainType,
			&job.BlockNumber,
			&job.Address,
			&job.Standard,
			&job.TokenID,
			&job.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve action job list: %w", err)
		}

		jobList = append(jobList, &job)
	}

	return jobList, nil
}

// ActionJob returns an Action job.
func (s *Store) ActionJob(id string) (*jobs.Action, error) {

	result, err := s.build.
		Select(actionJobsTableColumns...).
		From(actionJobsTableName).
		Where("id = ?", id).
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve action job: %w", ErrResourceNotFound)
	}

	var job jobs.Action
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.Standard,
		&job.TokenID,
		&job.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve action job: %w", err)
	}

	return &job, nil
}

// UpdateActionJobStatus updates an Action job status.
func (s *Store) UpdateActionJobStatus(id string, status jobs.Status) error {

	res, err := s.build.
		Update(actionJobsTableName).
		Where("id = ?", id).
		Set("status", status).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update action job state: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update action job state: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not update action job state: %w", ErrResourceNotFound)
	}

	return nil
}
