package api

import (
	"github.com/NFT-com/indexer/job"
)

type JobController interface {
	DiscoveryController
	ParsingController
}

type DiscoveryController interface {
	CreateDiscoveryJob(job.Discovery) (*job.Discovery, error)
	ListDiscoveryJobs(job.Status) ([]job.Discovery, error)
	GetDiscoveryJob(job.ID) (*job.Discovery, error)
	UpdateDiscoveryJobState(job.ID, job.Status) error
	RequeueDiscoveryJob(job.ID) (*job.Discovery, error)
}

type ParsingController interface {
	CreateParsingJob(job.Parsing) (*job.Parsing, error)
	ListParsingJobs(job.Status) ([]job.Parsing, error)
	GetParsingJob(job.ID) (*job.Parsing, error)
	UpdateParsingJobState(job.ID, job.Status) error
	RequeueParsingJob(job.ID) (*job.Parsing, error)
}
