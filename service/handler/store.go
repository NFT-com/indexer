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
	HighestBlockNumbersDiscoveryJob(chainURL, chainType string, addresses []string, standardType string) (map[string]string, error)
	DiscoveryJob(id string) (*jobs.Discovery, error)
	UpdateDiscoveryJobStatus(id string, status jobs.Status) error
}

// ParsingStore represents the parsing job store interface.
type ParsingStore interface {
	CreateParsingJob(job jobs.Parsing) error
	CreateParsingJobs(jobs []jobs.Parsing) error
	ParsingJobs(status jobs.Status) ([]jobs.Parsing, error)
	HighestBlockNumbersParsingJob(chainURL, chainType string, addresses []string, standardType, eventType string) (map[string]string, error)
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
