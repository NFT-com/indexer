package api

import (
	"github.com/NFT-com/indexer/job"
)

type JobController interface {
	DiscoveryController
	ParsingController
}

type DiscoveryController interface {
	CreateDiscoveryJob(job job.Discovery) (*job.Discovery, error)
	ListDiscoveryJobs(status job.Status) ([]job.Discovery, error)
	GetDiscoveryJob(id string) (*job.Discovery, error)
	UpdateDiscoveryJobState(id string, status job.Status) error
	RequeueDiscoveryJob(id string) (*job.Discovery, error)
}

type ParsingController interface {
	CreateParsingJob(job job.Parsing) (*job.Parsing, error)
	ListParsingJobs(status job.Status) ([]job.Parsing, error)
	GetParsingJob(id string) (*job.Parsing, error)
	UpdateParsingJobState(id string, status job.Status) error
	RequeueParsingJob(id string) (*job.Parsing, error)
}
