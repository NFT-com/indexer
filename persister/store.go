package persister

import (
	"github.com/NFT-com/indexer/jobs"
)

type Store interface {
	CreateParsingJobs(jobs []*jobs.Parsing) error
}
