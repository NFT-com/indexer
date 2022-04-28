package postgres

import (
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

// CreateParsingJob creates a new parsing job.
func (s *Store) CreateParsingJob(job *jobs.Parsing) error {

	_, err := s.build.
		Insert(parsingJobsTableName).
		Columns(parsingJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, job.Address, job.Standard, job.Event, job.Status).
		Exec()

	if err != nil {
		return fmt.Errorf("could not create parsing job: %w", err)
	}

	return nil
}

// CreateParsingJobs creates a batch of parsing jobs.
func (s *Store) CreateParsingJobs(jobs []*jobs.Parsing) error {

	query := s.build.
		Insert(parsingJobsTableName).
		Columns(parsingJobsTableColumns...)

	for _, job := range jobs {
		query = query.Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, job.Address, job.Standard, job.Event, job.Status)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not create parsing jobs: %w", err)
	}

	return nil
}

// ParsingJobs returns a list of parsing jobs filtered by status. Empty string status returns every job.
func (s *Store) ParsingJobs(status jobs.Status) ([]*jobs.Parsing, error) {

	query := s.build.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.OrderBy("block_number ASC")

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job list: %w", err)
	}
	defer result.Close()

	jobList := make([]*jobs.Parsing, 0)
	for result.Next() && result.Err() == nil {
		var job jobs.Parsing
		err = result.Scan(
			&job.ID,
			&job.ChainURL,
			&job.ChainID,
			&job.ChainType,
			&job.BlockNumber,
			&job.Address,
			&job.Standard,
			&job.Event,
			&job.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("could not retrieve parsing job list: %w", err)
		}

		jobList = append(jobList, &job)
	}

	return jobList, nil
}

// ParsingJob returns a parsing job.
func (s *Store) ParsingJob(id string) (*jobs.Parsing, error) {

	result, err := s.build.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName).
		Where("id = ?", id).
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
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.Standard,
		&job.Event,
		&job.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}

	return &job, nil
}

// HighestBlockNumberParsingJob returns the highest block number parsing job.
func (s *Store) HighestBlockNumberParsingJob(chainURL, chainType, address, Standard, eventType string) (*jobs.Parsing, error) {

	result, err := s.build.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName).
		Where("chain_url ILIKE ?", chainURL).
		Where("chain_type = ?", chainType).
		Where("address ILIKE ?", address).
		Where("interface_type = ?", Standard).
		Where("event_type ILIKE ?", eventType).
		OrderBy("block_number DESC").
		Limit(1).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve highest block number parsing job: %w", err)
	}
	defer result.Close()

	if !result.Next() {
		return nil, ErrResourceNotFound
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve highest block number parsing job: %w", result.Err())
	}

	var job jobs.Parsing
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.Standard,
		&job.Event,
		&job.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("could not scan highest block number parsing job: %w", err)
	}

	return &job, nil
}

// CountPendingParsingJobs returns the highest block number parsing job.
func (s *Store) CountPendingParsingJobs(chainURL, chainType, address, Standard, eventType string) (uint, error) {

	result, err := s.build.
		Select("COUNT(*)").
		From(parsingJobsTableName).
		Where("chain_url ILIKE ?", chainURL).
		Where("chain_type = ?", chainType).
		Where("address ILIKE ?", address).
		Where("interface_type = ?", Standard).
		Where("event_type ILIKE ?", eventType).
		Where("status IN (?, ?, ?)", jobs.StatusCreated, jobs.StatusProcessing, jobs.StatusQueued).
		Query()
	if err != nil {
		return 0, fmt.Errorf("could not count pending parsing jobs: %w", err)
	}
	defer result.Close()

	if !result.Next() {
		return 0, ErrResourceNotFound
	}

	if result.Err() != nil {
		return 0, fmt.Errorf("could not count pending parsing jobs: %w", result.Err())
	}

	var count uint
	err = result.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not scan pending parsing jobs: %w", err)
	}

	return count, nil
}

// UpdateParsingJobStatus updates a parsing job status.
func (s *Store) UpdateParsingJobStatus(id string, status jobs.Status) error {
	res, err := s.build.
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
		return fmt.Errorf("could not update parsing job status: %w", ErrResourceNotFound)
	}

	return nil
}
