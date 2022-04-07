package handler

import (
	"github.com/NFT-com/indexer/jobs"
)

// JobsStore represents the job store interface.
type JobsStore interface {
	DiscoveryStore
	ParsingStore
	AdditionStore
}

// DiscoveryStore represents the discovery job store interface.
type DiscoveryStore interface {
	CreateDiscoveryJob(job jobs.Discovery) error
	CreateDiscoveryJobs(jobs []jobs.Discovery) error
	DiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error)
	HighestBlockNumberDiscoveryJob(chainURL, chainType string, addresses []string, standardType, eventType string) (*jobs.Discovery, error)
	DiscoveryJob(id string) (*jobs.Discovery, error)
	UpdateDiscoveryJobStatus(id string, status jobs.Status) error
}

// ParsingStore represents the parsing job store interface.
type ParsingStore interface {
	CreateParsingJob(job jobs.Parsing) error
	CreateParsingJobs(jobs []jobs.Parsing) error
	ParsingJobs(status jobs.Status) ([]jobs.Parsing, error)
	HighestBlockNumberParsingJob(chainURL, chainType, address, standardType, eventType string) (*jobs.Parsing, error)
	ParsingJob(id string) (*jobs.Parsing, error)
	UpdateParsingJobStatus(id string, status jobs.Status) error
}

// AdditionStore represents the addition job store interface.
type AdditionStore interface {
	CreateAdditionJob(job jobs.Addition) error
	CreateAdditionJobs(jobs []jobs.Addition) error
	AdditionJobs(status jobs.Status) ([]jobs.Addition, error)
	AdditionJob(id string) (*jobs.Addition, error)
	UpdateAdditionJobStatus(id string, status jobs.Status) error
}
