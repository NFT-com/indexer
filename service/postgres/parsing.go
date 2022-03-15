package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/job"
)

func (s *Store) CreateParsingJob(job job.Parsing) error {
	_, err := s.sqlBuilder.
		Insert(ParsingJobsDBName).
		Columns(ParsingJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, job.Address, job.StandardType, job.EventType, job.Status).
		Exec()

	if err != nil {
		return fmt.Errorf("failed to create parsing job: %v", err)
	}

	return nil
}

func (s *Store) ParsingJobs(status job.Status) ([]job.Parsing, error) {
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

	jobList := make([]job.Parsing, 0)
	for result.Next() && result.Err() == nil {
		var job job.Parsing
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

func (s *Store) ParsingJob(jobID job.ID) (*job.Parsing, error) {
	result, err := s.sqlBuilder.
		Select(ParsingJobsTableColumns...).
		From(ParsingJobsDBName).
		Where("id = ?", jobID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parsing job: %v", err)
	}

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("failed to retrieve parsing job: %v", ErrResourceNotFound)
	}

	var job job.Parsing
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
		return nil, fmt.Errorf("failed to retrieve parsing job: %v", err)
	}

	return &job, nil
}

func (s *Store) UpdateParsingJobState(jobID job.ID, jobStatus job.Status) error {
	res, err := s.sqlBuilder.
		Update(ParsingJobsDBName).
		Where("id = ?", jobID).
		Set("status", jobStatus).
		Set("updated_at", time.Now()).
		Exec()

	if err != nil {
		return fmt.Errorf("failed to update parsing job state: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update parsing job state: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed to update parsing job state: %v", ErrResourceNotFound)
	}

	return nil
}
