package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingRepository interface {
	Find(wheres ...string) ([]*jobs.Parsing, error)
	UpdateStatus(status string, parsingIDs ...string) error
}
