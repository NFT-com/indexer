package controller

import (
	"github.com/NFT-com/indexer/job"
)

type JobsStore interface {
	CreateDiscoveryJob(job.Discovery) error
	ListDiscoveryJobs(job.Status) ([]job.Discovery, error)
	GetDiscoveryJob(job.ID) (*job.Discovery, error)
	UpdateDiscoveryJobState(job.ID, job.Status) error

	CreateParsingJob(job.Parsing) error
	ListParsingJobs(job.Status) ([]job.Parsing, error)
	GetParsingJob(job.ID) (*job.Parsing, error)
	UpdateParsingJobState(job.ID, job.Status) error
}
