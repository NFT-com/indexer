package creator

import (
	"github.com/NFT-com/indexer/jobs"
)

type Checker interface {
	CountPendingParsingJobs(chainURL, chainType, address, Standard, eventType string) (uint, error)
	HighestBlockNumberParsingJob(chainURL, chainType, address, Standard, eventType string) (*jobs.Parsing, error)
}
