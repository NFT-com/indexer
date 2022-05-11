package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingStore interface {
	List(status string) ([]*jobs.Parsing, error)
	UpdateStatus(status string, statusMessage string, parsingIDs ...string) error
}

type ActionStore interface {
	List(status string) ([]*jobs.Action, error)
	UpdateStatus(status string, statusMessage string, actionIDs ...string) error
}
