package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ActionRepository interface {
	Find(wheres ...string) ([]*jobs.Action, error)
	UpdateStatus(status string, actionIDs ...string) error
}
