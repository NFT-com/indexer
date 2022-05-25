package pipeline

import (
	"github.com/NFT-com/indexer/models/jobs"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type ParsingStore interface {
	Update(selector storage.UpdateSelector, setters ...storage.UpdateSetter) error
}

type ActionStore interface {
	Insert(actions ...*jobs.Action) error
	Update(selector storage.UpdateSelector, setters ...storage.UpdateSetter) error
}
