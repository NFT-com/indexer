package creator

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type JobRepository interface {
	Count(wheres ...string) (uint, error)
	Top(order string, wheres ...string) (*jobs.Parsing, error)
}
