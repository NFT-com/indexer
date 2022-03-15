package controller

import (
	"github.com/NFT-com/indexer/jobs"
)

type JobsStore interface {
	DiscoveryStore
	ParsingStore
}

type DiscoveryStore interface {
	CreateDiscoveryJob(job jobs.Discovery) error
	DiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error)
	DiscoveryJob(id string) (*jobs.Discovery, error)
	UpdateDiscoveryJobState(id string, status jobs.Status) error
}

type ParsingStore interface {
	CreateParsingJob(job jobs.Parsing) error
	ParsingJobs(status jobs.Status) ([]jobs.Parsing, error)
	ParsingJob(id string) (*jobs.Parsing, error)
	UpdateParsingJobState(id string, status jobs.Status) error
}
