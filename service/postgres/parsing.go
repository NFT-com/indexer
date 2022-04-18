package postgres

import (
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/NFT-com/indexer/jobs"
)

// CreateParsingJob creates a new parsing job.
func (s *Store) CreateParsingJob(job jobs.Parsing) error {
	_, err := s.sqlBuilder.
		Insert(parsingJobsTableName).
		Columns(parsingJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, pq.Array(job.Addresses), job.StandardType, job.EventType, job.Status).
		Exec()

	if err != nil {
		return fmt.Errorf("could not create parsing job: %w", err)
	}

	return nil
}

// CreateParsingJobs creates a batch of parsing jobs.
func (s *Store) CreateParsingJobs(jobs []jobs.Parsing) error {
	query := s.sqlBuilder.
		Insert(parsingJobsTableName).
		Columns(parsingJobsTableColumns...)

	for _, job := range jobs {
		query = query.Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, pq.Array(job.Addresses), job.StandardType, job.EventType, job.Status)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not create parsing jobs: %w", err)
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
			&job.ChainID,
			&job.ChainType,
			&job.BlockNumber,
			pq.Array(&job.Addresses),
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
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		pq.Array(&job.Addresses),
		&job.StandardType,
		&job.EventType,
		&job.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}

	return &job, nil
}

// HighestBlockNumbersParsingJob returns the highest block numbers for parsing jobs.
func (s *Store) HighestBlockNumbersParsingJob(chainURL, chainType string, addresses []string, standardType, eventType string) (map[string]string, error) {
	result, err := s.sqlBuilder.
		Select("jobs.address, jobs.block_number").
		FromSelect(
			s.sqlBuilder.
				Select("unnest(addresses) as address, max(block_number) as block_number").
				From(parsingJobsTableName).
				Where("chain_url = ?", chainURL).
				Where("chain_type = ?", chainType).
				Where("? <@ addresses", pq.Array(addresses)).
				Where("interface_type = ?", standardType).
				GroupBy("address"),
			"jobs",
		).
		Where("jobs.address = any(?)", pq.Array(addresses)).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve highest block number parsing job: %w", err)
	}
	defer result.Close()

	highestBlocks := make(map[string]string)
	for result.Next() && result.Err() == nil {
		var address string
		var value string
		err = result.Scan(
			&address,
			&value,
		)

		if err != nil {
			return nil, fmt.Errorf("could not heighest block: %w", err)
		}

		highestBlocks[address] = value
	}

	if err != nil {
		return nil, fmt.Errorf("could not retrieve highest block number for addresses: %w", err)
	}

	return highestBlocks, nil
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
