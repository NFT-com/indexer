package creator

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type Persister interface {
	Store(job *jobs.Parsing)
}
