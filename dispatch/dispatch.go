package dispatch

import "github.com/NFT-com/indexer/event"

type Dispatcher interface {
	Dispatch(function string, event *event.Event) error
}
