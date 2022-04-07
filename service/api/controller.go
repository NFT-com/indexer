package api

import (
	"github.com/NFT-com/indexer/jobs"
)

// JobsHandler represents the business layer of the jobs api.
type JobsHandler interface {
	DiscoveryJobHandler
	ParsingJobHandler
}

// DiscoveryJobHandler represent the business layer of the discovery jobs api.
type DiscoveryJobHandler interface {
	CreateDiscoveryJob(job jobs.Discovery) (*jobs.Discovery, error)
	ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error)
	GetDiscoveryJob(id string) (*jobs.Discovery, error)
	GetHighestBlockNumberDiscoveryJob(chainURL, chainType string, addresses []string, standardType, eventType string) (*jobs.Discovery, error)
	UpdateDiscoveryJobStatus(id string, status jobs.Status) error
}

// ParsingJobHandler represent the business layer of the parsing jobs api.
type ParsingJobHandler interface {
	CreateParsingJob(job jobs.Parsing) (*jobs.Parsing, error)
	ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error)
	GetParsingJob(id string) (*jobs.Parsing, error)
	GetHighestBlockNumberParsingJob(chainURL, chainType, address, standardType, eventType string) (*jobs.Parsing, error)
	UpdateParsingJobStatus(id string, status jobs.Status) error
}
