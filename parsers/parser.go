package parsers

import (
	"github.com/NFT-com/indexer/events"
)

type Parser interface {
	ParseRawEvent(rawEvent events.RawEvent) (events.Event, error)
}
