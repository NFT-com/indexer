package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/NFT-com/indexer/jobs"
)

func (s *Store) CreateDiscoveryJob(job jobs.Discovery) error {
	_, err := s.sqlBuilder.
		Insert(discoveryJobsTableName).
		Columns(discoveryJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, pq.Array(job.Addresses), job.StandardType, job.Status).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create discovery job: %w", err)
	}

	return nil
}

func (s *Store) DiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error) {
	query := s.sqlBuilder.
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

	jobList := make([]jobs.Discovery, 0)
	for result.Next() && result.Err() == nil {
		var job jobs.Discovery

		err = result.Scan(
			&job.ID,
			&job.ChainURL,
			&job.ChainType,
			&job.BlockNumber,
			pq.Array(&job.Addresses),
			&job.StandardType,
			&job.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve discovery job list: %w", err)
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

func (s *Store) DiscoveryJob(id string) (*jobs.Discovery, error) {
	result, err := s.sqlBuilder.
		Select(discoveryJobsTableColumns...).
		From(discoveryJobsTableName).
		Where("id = ?", id).
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %w", errResourceNotFound)
	}

	var job jobs.Discovery
	rawAddresses := make([]byte, 0)

	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainType,
		&job.BlockNumber,
		&rawAddresses,
		&job.StandardType,
		&job.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %w", err)
	}

	err = json.Unmarshal(rawAddresses, &job.Addresses)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %w", err)
	}

	return &job, nil
}

func (s *Store) UpdateDiscoveryJobState(id string, status jobs.Status) error {
	res, err := s.sqlBuilder.
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
		return fmt.Errorf("could not update discovery job state: %w", errResourceNotFound)
	}

	return nil
}
