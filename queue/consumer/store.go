package consumer

import (
	"github.com/NFT-com/indexer/events"
)

type Store interface {
	InsertRawEvent(event events.RawEvent) error
	InsertNewNFT(chain, contract, id, owner string) error
	UpdateNFT(chain, contract, id, owner string) error
}
