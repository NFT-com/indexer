package pipeline

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingStore interface {
	UpdateStatus(parsingIDs []string, status string) error
}

type ActionStore interface {
	Insert(action *jobs.Action) error
}
