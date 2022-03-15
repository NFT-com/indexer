package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

func (s *Store) CreateParsingJob(job jobs.Parsing) error {
	_, err := s.sqlBuilder.
		Insert(ParsingJobsDBName).
		Columns(ParsingJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, job.Address, job.StandardType, job.EventType, job.Status).
		Exec()

	if err != nil {
		return fmt.Errorf("failed to create parsing jobs: %v", err)
	}

	return nil
}

func (s *Store) ParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	query := s.sqlBuilder.
		Select(ParsingJobsTableColumns...).
		From(ParsingJobsDBName)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parsing jobs list: %v", err)
	}

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
			return nil, fmt.Errorf("failed to retrieve parsing jobs list: %v", err)
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

func (s *Store) ParsingJob(jobID string) (*jobs.Parsing, error) {
	result, err := s.sqlBuilder.
		Select(ParsingJobsTableColumns...).
		From(ParsingJobsDBName).
		Where("id = ?", jobID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parsing jobs: %v", err)
	}

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("failed to retrieve parsing jobs: %v", ErrResourceNotFound)
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
		return nil, fmt.Errorf("failed to retrieve parsing jobs: %v", err)
	}

	return &job, nil
}

func (s *Store) UpdateParsingJobState(jobID string, jobStatus jobs.Status) error {
	res, err := s.sqlBuilder.
		Update(ParsingJobsDBName).
		Where("id = ?", jobID).
		Set("status", jobStatus).
		Set("updated_at", time.Now()).
		Exec()

	if err != nil {
		return fmt.Errorf("failed to update parsing jobs state: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update parsing jobs state: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed to update parsing jobs state: %v", ErrResourceNotFound)
	}

	return nil
}
