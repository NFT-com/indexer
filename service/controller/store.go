package controller

import (
	"github.com/NFT-com/indexer/job"
)

type DiscoveryJobsStore interface {
	CreateDiscoveryJob(job.Discovery) error
	ListDiscoveryJobs(job.Status) ([]job.Discovery, error)
	GetDiscoveryJob(job.ID) (job.Discovery, error)
	UpdateDiscoveryJobState(job.ID, job.Status) error
}

type ParsingJobsStore interface {
	CreateParsingJob(job.Parsing) error
	ListParsingJobs(job.Status) ([]job.Parsing, error)
	GetParsingJob(job.ID) (job.Parsing, error)
	UpdateParsingJobState(job.ID, job.Status) error
}
