package consumer

import (
	"github.com/NFT-com/indexer/event"
)

type Store interface {
	InsertHistory(event event.Event) error
}
