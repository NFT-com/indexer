package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type ParsingStore interface {
	List(status string) ([]*jobs.Parsing, error)
	UpdateStatus(status string, parsingIDs []string, options ...storage.UpdateStatusOption) error
}

type ActionStore interface {
	List(status string) ([]*jobs.Action, error)
	UpdateStatus(status string, actionIDs []string, options ...storage.UpdateStatusOption) error
}
