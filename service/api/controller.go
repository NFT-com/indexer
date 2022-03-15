package api

import (
	"github.com/NFT-com/indexer/jobs"
)

type JobController interface {
	DiscoveryController
	ParsingController
}

type DiscoveryController interface {
	CreateDiscoveryJob(jobs.Discovery) (*jobs.Discovery, error)
	ListDiscoveryJobs(jobs.Status) ([]jobs.Discovery, error)
	GetDiscoveryJob(jobs.ID) (*jobs.Discovery, error)
	UpdateDiscoveryJobState(jobs.ID, jobs.Status) error
	RequeueDiscoveryJob(jobs.ID) (*jobs.Discovery, error)
}

type ParsingController interface {
	CreateParsingJob(jobs.Parsing) (*jobs.Parsing, error)
	ListParsingJobs(jobs.Status) ([]jobs.Parsing, error)
	GetParsingJob(jobs.ID) (*jobs.Parsing, error)
	UpdateParsingJobState(jobs.ID, jobs.Status) error
	RequeueParsingJob(jobs.ID) (*jobs.Parsing, error)
}
