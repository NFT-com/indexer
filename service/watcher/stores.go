package watcher

import (
	"github.com/NFT-com/indexer/models/jobs"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type ParsingStore interface {
	List(status string) ([]*jobs.Parsing, error)
	Update(selector storage.UpdateSelector, setters ...storage.UpdateSetter) error
}

type ActionStore interface {
	List(status string) ([]*jobs.Action, error)
	Update(selector storage.UpdateSelector, setters ...storage.UpdateSetter) error
}
