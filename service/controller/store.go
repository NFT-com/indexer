package controller

import (
	"github.com/NFT-com/indexer/job"
)

type JobsStore interface {
	DiscoveryStore
	ParsingStore
}

type DiscoveryStore interface {
	CreateDiscoveryJob(job.Discovery) error
	DiscoveryJobs(job.Status) ([]job.Discovery, error)
	DiscoveryJob(job.ID) (*job.Discovery, error)
	UpdateDiscoveryJobState(job.ID, job.Status) error
}

type ParsingStore interface {
	CreateParsingJob(job.Parsing) error
	ParsingJobs(job.Status) ([]job.Parsing, error)
	ParsingJob(job.ID) (*job.Parsing, error)
	UpdateParsingJobState(job.ID, job.Status) error
}
