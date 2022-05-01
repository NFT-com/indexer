package pipeline

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingStore interface {
	Update(parsing *jobs.Parsing) error
}
