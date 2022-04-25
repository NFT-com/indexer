package creator

import (
	"github.com/NFT-com/indexer/jobs"
)

type Checker interface {
	CountPendingParsingJobs(chainURL, chainType, address, Standard, eventType string) (uint, error)
	LastParsingJob(chainID string, address string, eventType string) (*jobs.Parsing, error)
}
