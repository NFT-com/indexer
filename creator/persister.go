package creator

import (
	"github.com/NFT-com/indexer/jobs"
)

type Persister interface {
	Store(job *jobs.Parsing)
}
