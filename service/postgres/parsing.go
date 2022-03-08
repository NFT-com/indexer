package postgres

import (
	"time"

	"github.com/NFT-com/indexer/job"
)

func (s *Store) CreateParsingJob(job job.Parsing) error {
	_, err := s.sqlBuilder.
		Insert(ParsingJobsDBName).
		Columns(ParsingJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, job.Address, job.InterfaceType, job.EventType, job.Status).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ListParsingJobs(status job.Status) ([]job.Parsing, error) {
	query := s.sqlBuilder.
		Select(ParsingJobsTableColumns...).
		From(ParsingJobsDBName)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, err
	}

	jobList := make([]job.Parsing, 0)
	for result.Next() && result.Err() == nil {
		parsingJob := job.Parsing{}

		err := result.Scan(
			&parsingJob.ID,
			&parsingJob.ChainURL,
			&parsingJob.ChainType,
			&parsingJob.BlockNumber,
			&parsingJob.Address,
			&parsingJob.InterfaceType,
			&parsingJob.EventType,
			&parsingJob.Status,
		)

		if err != nil {
			return nil, err
		}

		jobList = append(jobList, parsingJob)
	}

	return jobList, nil
}

func (s *Store) GetParsingJob(jobID job.ID) (job.Parsing, error) {
	result, err := s.sqlBuilder.
		Select(ParsingJobsTableColumns...).
		From(ParsingJobsDBName).
		Where("id = ?", jobID).
		Query()
	if err != nil {
		return job.Parsing{}, err
	}

	if !result.Next() || result.Err() != nil {
		return job.Parsing{}, ErrResourceNotFound
	}

	parsingJob := job.Parsing{}

	err = result.Scan(
		&parsingJob.ID,
		&parsingJob.ChainURL,
		&parsingJob.ChainType,
		&parsingJob.BlockNumber,
		&parsingJob.Address,
		&parsingJob.InterfaceType,
		&parsingJob.EventType,
		&parsingJob.Status,
	)

	if err != nil {
		return job.Parsing{}, err
	}

	return parsingJob, nil
}

func (s *Store) UpdateParsingJobState(jobID job.ID, jobStatus job.Status) error {
	_, err := s.sqlBuilder.
		Update(ParsingJobsDBName).
		Where("id = ?", jobID).
		Set("status", jobStatus).
		Set("updated_at", time.Now()).
		Exec()

	if err != nil {
		return err
	}

	return nil
}