package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type ActionRepository interface {
	Find(wheres ...string) ([]*jobs.Action, error)
	UpdateStatus(actionID string, status string) error
}
