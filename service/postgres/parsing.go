package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

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

func (s *Store) ParsingJob(jobID jobs.ID) (*jobs.Parsing, error) {
	result, err := s.sqlBuilder.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName).
		Where("id = ?", jobID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", ErrResourceNotFound)
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

func (s *Store) UpdateParsingJobState(jobID jobs.ID, jobStatus jobs.Status) error {
	res, err := s.sqlBuilder.
		Update(parsingJobsTableName).
		Where("id = ?", jobID).
		Set("status", jobStatus).
		Set("updated_at", time.Now()).
		Exec()

	if err != nil {
		return fmt.Errorf("could not update parsing job state: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update parsing job state: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not update parsing job state: %w", ErrResourceNotFound)
	}

	return nil
}
