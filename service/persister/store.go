package persister

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type Store interface {
	CreateParsingJobs(jobs []*jobs.Parsing) error
}
