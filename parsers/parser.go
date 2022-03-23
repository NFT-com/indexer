package parsers

import (
	"github.com/NFT-com/indexer/event"
)

type Parser interface {
	ParseRawEvent(rawEvent event.RawEvent) (*event.Event, error)
}
