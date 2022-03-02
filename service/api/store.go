package api

import (
	"github.com/NFT-com/indexer/job"
)

type DiscoveryJobsStore interface {
	CreateDiscoveryJob(job.Discovery) error
	ListDiscoveryJobs(job.Status) ([]job.Discovery, error)
	GetDiscoveryJob(job.ID) (job.Discovery, error)
	CancelDeliveryJob(job.ID) error
}

type ParsingJobsStore interface {
	CreateParsingJob(job.Parsing) error
	ListParsingJobs(job.Status) ([]job.Parsing, error)
	GetParsingJob(job.ID) (job.Parsing, error)
	CancelParsingJob(job.ID) error
}
