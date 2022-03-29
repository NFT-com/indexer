package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

// CreateParsingJob creates a new parsing job.
func (s *Store) CreateParsingJob(job jobs.Parsing) error {
	_, err := s.sqlBuilder.
		Insert(parsingJobsTableName).
		Columns(parsingJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, job.Address, job.StandardType, job.EventType, job.Status).
		Exec()

	if err != nil {
		return fmt.Errorf("could not create parsing job: %w", err)
	}

	return nil
}

// ParsingJobs returns a list of parsing jobs filtered by status. Empty string status returns every job.
func (s *Store) ParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	query := s.sqlBuilder.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job list: %w", err)
	}
	defer result.Close()

	jobList := make([]jobs.Parsing, 0)
	for result.Next() && result.Err() == nil {
		var job jobs.Parsing
		err = result.Scan(
			&job.ID,
			&job.ChainURL,
			&job.ChainType,
			&job.BlockNumber,
			&job.Address,
			&job.StandardType,
			&job.EventType,
			&job.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("could not retrieve parsing job list: %w", err)
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

// ParsingJob returns a parsing job.
func (s *Store) ParsingJob(id string) (*jobs.Parsing, error) {
	result, err := s.sqlBuilder.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName).
		Where("id = ?", id).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", errResourceNotFound)
	}

	var job jobs.Parsing
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.StandardType,
		&job.EventType,
		&job.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}

	return &job, nil
}

// UpdateParsingJobStatus updates a parsing job status.
func (s *Store) UpdateParsingJobStatus(id string, status jobs.Status) error {
	res, err := s.sqlBuilder.
		Update(parsingJobsTableName).
		Where("id = ?", id).
		Set("status", status).
		Set("updated_at", time.Now()).
		Exec()

	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not update parsing job status: %w", errResourceNotFound)
	}

	return nil
}
