package pipeline

import (
	"github.com/NFT-com/indexer/models/jobs"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type ParsingStore interface {
	UpdateStatus(status string, parsingIDs []string, options ...storage.UpdateStatusOption) error
}

type ActionStore interface {
	Insert(actions ...*jobs.Action) error
	UpdateStatus(status string, actionIDs []string, options ...storage.UpdateStatusOption) error
}
