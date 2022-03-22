package api

import (
	"github.com/NFT-com/indexer/jobs"
)

type JobsHandler interface {
	DiscoveryJobHandler
	ParsingJobHandler
}

type DiscoveryJobHandler interface {
	CreateDiscoveryJob(jobs.Discovery) (*jobs.Discovery, error)
	ListDiscoveryJobs(jobs.Status) ([]jobs.Discovery, error)
	GetDiscoveryJob(jobs.ID) (*jobs.Discovery, error)
	UpdateDiscoveryJobState(jobs.ID, jobs.Status) error
	RequeueDiscoveryJob(jobs.ID) (*jobs.Discovery, error)
}

type ParsingJobHandler interface {
	CreateParsingJob(jobs.Parsing) (*jobs.Parsing, error)
	ListParsingJobs(jobs.Status) ([]jobs.Parsing, error)
	GetParsingJob(jobs.ID) (*jobs.Parsing, error)
	UpdateParsingJobState(jobs.ID, jobs.Status) error
	RequeueParsingJob(jobs.ID) (*jobs.Parsing, error)
}
