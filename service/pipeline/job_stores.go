package pipeline

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingStore interface {
	UpdateStatus(status string, statusMessage string, parsingIDs ...string) error
}

type ActionStore interface {
	Insert(actions ...*jobs.Action) error
	UpdateStatus(status string, statusMessage string, actionIDs ...string) error
}
