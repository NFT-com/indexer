package api

import (
	"github.com/NFT-com/indexer/jobs"
)

type JobsHandler interface {
	DiscoveryJobHandler
	ParsingJobHandler
}

type DiscoveryJobHandler interface {
	CreateDiscoveryJob(job jobs.Discovery) (*jobs.Discovery, error)
	ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error)
	GetDiscoveryJob(id string) (*jobs.Discovery, error)
	UpdateDiscoveryJobState(id string, status jobs.Status) error
	RequeueDiscoveryJob(id string) (*jobs.Discovery, error)
}

type ParsingJobHandler interface {
	CreateParsingJob(job jobs.Parsing) (*jobs.Parsing, error)
	ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error)
	GetParsingJob(id string) (*jobs.Parsing, error)
	UpdateParsingJobState(id string, status jobs.Status) error
	RequeueParsingJob(id string) (*jobs.Parsing, error)
}
