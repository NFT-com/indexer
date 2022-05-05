package persister

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingStore interface {
	Insert(parsings ...*jobs.Parsing) error
	UpdateStatus(status string, parsingIDs ...string) error
}

type ActionStore interface {
	Insert(actions ...*jobs.Action) error
	UpdateStatus(status string, actionIDs ...string) error
}
