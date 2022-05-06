package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingStore interface {
	Find(wheres ...string) ([]*jobs.Parsing, error)
	UpdateStatus(status string, parsingIDs ...string) error
}

type ActionStore interface {
	Find(wheres ...string) ([]*jobs.Action, error)
	UpdateStatus(status string, actionIDs ...string) error
}
