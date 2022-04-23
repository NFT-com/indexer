package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

// CreateAdditionJob creates a new addition job.
func (s *Store) CreateAdditionJob(job jobs.Addition) error {

	_, err := s.build.
		Insert(additionJobsTableName).
		Columns(additionJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, job.Address, job.StandardType, job.TokenID, job.Status).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create addition job: %w", err)
	}

	return nil
}

// CreateAdditionJobs creates a batch of addition jobs.
func (s *Store) CreateAdditionJobs(jobs []jobs.Addition) error {

	query := s.build.
		Insert(additionJobsTableName).
		Columns(additionJobsTableColumns...)

	for _, job := range jobs {
		query = query.Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, job.Address, job.StandardType, job.TokenID, job.Status)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not create addition jobs: %w", err)
	}

	return nil
}

// AdditionJobs returns a list of addition jobs filtered by status. Empty string status returns every job.
func (s *Store) AdditionJobs(status jobs.Status) ([]jobs.Addition, error) {

	query := s.build.
		Select(additionJobsTableColumns...).
		From(additionJobsTableName)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve addition job list: %w", err)
	}
	defer result.Close()

	jobList := make([]jobs.Addition, 0)
	for result.Next() && result.Err() == nil {
		var job jobs.Addition
		err = result.Scan(
			&job.ID,
			&job.ChainURL,
			&job.ChainID,
			&job.ChainType,
			&job.BlockNumber,
			&job.Address,
			&job.StandardType,
			&job.TokenID,
			&job.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve addition job list: %w", err)
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

// AdditionJob returns an addition job.
func (s *Store) AdditionJob(id string) (*jobs.Addition, error) {

	result, err := s.build.
		Select(additionJobsTableColumns...).
		From(additionJobsTableName).
		Where("id = ?", id).
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve addition job: %w", errResourceNotFound)
	}

	var job jobs.Addition
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.StandardType,
		&job.TokenID,
		&job.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve addition job: %w", err)
	}

	return &job, nil
}

// UpdateAdditionJobStatus updates an addition job status.
func (s *Store) UpdateAdditionJobStatus(id string, status jobs.Status) error {

	res, err := s.build.
		Update(additionJobsTableName).
		Where("id = ?", id).
		Set("status", status).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update addition job state: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update addition job state: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not update addition job state: %w", errResourceNotFound)
	}

	return nil
}
