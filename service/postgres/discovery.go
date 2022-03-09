package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NFT-com/indexer/job"
)

func (s *Store) CreateDiscoveryJob(job job.Discovery) error {
	rawAddresses, err := json.Marshal(job.Addresses)
	if err != nil {
		return err
	}

	_, err = s.sqlBuilder.
		Insert(DiscoveryJobsDBName).
		Columns(DiscoveryJobsTableColumns...).
		Values(job.ID, job.ChainURL, job.ChainType, job.BlockNumber, rawAddresses, job.StandardType, job.Status).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to create discovery job: %v", err)
	}

	return nil
}

func (s *Store) ListDiscoveryJobs(status job.Status) ([]job.Discovery, error) {
	query := s.sqlBuilder.
		Select(DiscoveryJobsTableColumns...).
		From(DiscoveryJobsDBName)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discovery job list: %v", err)
	}

	jobList := make([]job.Discovery, 0)
	for result.Next() && result.Err() == nil {
		discoveryJob := job.Discovery{}
		rawAddresses := make([]byte, 0)

		err = result.Scan(
			&discoveryJob.ID,
			&discoveryJob.ChainURL,
			&discoveryJob.ChainType,
			&discoveryJob.BlockNumber,
			&rawAddresses,
			&discoveryJob.StandardType,
			&discoveryJob.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve discovery job list: %v", err)
		}

		err = json.Unmarshal(rawAddresses, &discoveryJob.Addresses)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve discovery job list: %v", err)
		}

		jobList = append(jobList, discoveryJob)
	}

	return jobList, nil
}

func (s *Store) GetDiscoveryJob(jobID job.ID) (*job.Discovery, error) {
	result, err := s.sqlBuilder.
		Select(DiscoveryJobsTableColumns...).
		From(DiscoveryJobsDBName).
		Where("id = ?", jobID).
		Query()
	if err != nil {
		return nil, err
	}

	if !result.Next() || result.Err() != nil {
		return nil, fmt.Errorf("failed to retrieve discovery job: %v", ErrResourceNotFound)
	}

	discoveryJob := job.Discovery{}
	rawAddresses := make([]byte, 0)

	err = result.Scan(
		&discoveryJob.ID,
		&discoveryJob.ChainURL,
		&discoveryJob.ChainType,
		&discoveryJob.BlockNumber,
		&rawAddresses,
		&discoveryJob.StandardType,
		&discoveryJob.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discovery job: %v", err)
	}

	err = json.Unmarshal(rawAddresses, &discoveryJob.Addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discovery job: %v", err)
	}

	return &discoveryJob, nil
}

func (s *Store) UpdateDiscoveryJobState(jobID job.ID, jobStatus job.Status) error {
	_, err := s.sqlBuilder.
		Update(DiscoveryJobsDBName).
		Where("id = ?", jobID).
		Set("status", jobStatus).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to update discovery job state: %v", err)
	}

	return nil
}
