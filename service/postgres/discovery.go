package postgres

import (
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/NFT-com/indexer/jobs"
)

// CreateDiscoveryJob creates a new discovery job.
func (s *Store) CreateDiscoveryJob(job *jobs.Discovery) error {

	_, err := s.build.
		Insert(discoveryJobsTableName).
		Columns(discoveryJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, pq.Array(job.Addresses), job.StandardType, job.Status).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create discovery job: %w", err)
	}

	return nil
}

// CreateDiscoveryJobs creates a batch of discovery jobs.
func (s *Store) CreateDiscoveryJobs(jobs []*jobs.Discovery) error {

	query := s.build.
		Insert(discoveryJobsTableName).
		Columns(discoveryJobsTableColumns...)

	for _, job := range jobs {
		query = query.Values(job.ID, job.ChainURL, job.ChainID, job.ChainType, job.BlockNumber, pq.Array(job.Addresses), job.StandardType, job.Status)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not create discovery jobs: %w", err)
	}

	return nil
}

// DiscoveryJobs returns a list of discovery jobs filtered by status. Empty string status returns every job.
func (s *Store) DiscoveryJobs(status jobs.Status) ([]*jobs.Discovery, error) {

	query := s.build.
		Select(discoveryJobsTableColumns...).
		From(discoveryJobsTableName)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve discovery job list: %w", err)
	}
	defer result.Close()

	jobList := make([]*jobs.Discovery, 0)
	for result.Next() && result.Err() == nil {
		var job jobs.Discovery

		err = result.Scan(
			&job.ID,
			&job.ChainURL,
			&job.ChainID,
			&job.ChainType,
			&job.BlockNumber,
			pq.Array(&job.Addresses),
			&job.StandardType,
			&job.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve discovery job list: %w", err)
		}

		jobList = append(jobList, &job)
	}

	return jobList, nil
}

// DiscoveryJob returns a discovery job.
func (s *Store) DiscoveryJob(id string) (*jobs.Discovery, error) {

	result, err := s.build.
		Select(discoveryJobsTableColumns...).
		From(discoveryJobsTableName).
		Where("id = ?", id).
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %w", ErrResourceNotFound)
	}

	var job jobs.Discovery
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		pq.Array(&job.Addresses),
		&job.StandardType,
		&job.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %w", err)
	}

	return &job, nil
}

// HighestBlockNumberDiscoveryJob returns the highest block number discovery job.
func (s *Store) HighestBlockNumberDiscoveryJob(chainURL, chainType string, addresses []string, standardType, eventType string) (*jobs.Discovery, error) {

	result, err := s.build.
		Select(parsingJobsTableColumns...).
		From(parsingJobsTableName).
		Where("chain_url = ?", chainURL).
		Where("chain_type = ?", chainType).
		Where("addresses <@ ? AND ? <@ addresses", pq.Array(addresses), pq.Array(addresses)).
		Where("standard_type = ?", standardType).
		Where("event_type ILIKE ?", eventType).
		OrderBy("block_number DESC").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve highest block number discovery job: %w", err)
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve highest block number discovery job: %w", ErrResourceNotFound)
	}

	var job jobs.Discovery
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		pq.Array(&job.Addresses),
		&job.StandardType,
		&job.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve highest block number discovery job: %w", err)
	}

	return &job, nil
}

// UpdateDiscoveryJobStatus updates a discovery job status.
func (s *Store) UpdateDiscoveryJobStatus(id string, status jobs.Status) error {

	res, err := s.build.
		Update(discoveryJobsTableName).
		Where("id = ?", id).
		Set("status", status).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update discovery job state: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update discovery job state: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not update discovery job state: %w", ErrResourceNotFound)
	}

	return nil
}
