package handler

import (
	"github.com/NFT-com/indexer/jobs"
)

type JobsStore interface {
	DiscoveryStore
	ParsingStore
}

type DiscoveryStore interface {
	CreateDiscoveryJob(jobs.Discovery) error
	DiscoveryJobs(jobs.Status) ([]jobs.Discovery, error)
	DiscoveryJob(jobs.ID) (*jobs.Discovery, error)
	UpdateDiscoveryJobState(jobs.ID, jobs.Status) error
}

type ParsingStore interface {
	CreateParsingJob(jobs.Parsing) error
	ParsingJobs(jobs.Status) ([]jobs.Parsing, error)
	ParsingJob(jobs.ID) (*jobs.Parsing, error)
	UpdateParsingJobState(jobs.ID, jobs.Status) error
}
