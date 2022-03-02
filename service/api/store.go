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
}
