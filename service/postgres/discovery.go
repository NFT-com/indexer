package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NFT-com/indexer/jobs"
	"github.com/lib/pq"
)

func (s *Store) CreateDiscoveryJob(job jobs.Discovery) error {
	_, err := s.sqlBuilder.
		Insert(discoveryJobsTableName).
		Columns(discoveryJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, pq.Array(job.Addresses), job.StandardType, job.Status).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create discovery job: %v", err)
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
		return nil, fmt.Errorf("could not retrieve discovery job list: %v", err)
	}

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
			return nil, fmt.Errorf("could not retrieve discovery job list: %v", err)
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

func (s *Store) DiscoveryJob(jobID jobs.ID) (*jobs.Discovery, error) {
	result, err := s.sqlBuilder.
		Select(discoveryJobsTableColumns...).
		From(discoveryJobsTableName).
		Where("id = ?", jobID).
		Query()
	if err != nil {
		return nil, err
	}

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %v", ErrResourceNotFound)
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
		return nil, fmt.Errorf("could not retrieve discovery job: %v", err)
	}

	err = json.Unmarshal(rawAddresses, &job.Addresses)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve discovery job: %v", err)
	}

	return &job, nil
}

func (s *Store) UpdateDiscoveryJobState(jobID jobs.ID, jobStatus jobs.Status) error {
	res, err := s.sqlBuilder.
		Update(discoveryJobsTableName).
		Where("id = ?", jobID).
		Set("status", jobStatus).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update discovery job state: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update discovery job state: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not update discovery job state: %v", ErrResourceNotFound)
	}

	return nil
}
