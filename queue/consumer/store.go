package consumer

import (
	"github.com/NFT-com/indexer/events"
)

type Store interface {
	InsertRawEvent(event events.RawEvent) error
	InsertNewNFT(network, chain, contract, id, owner string) error
	UpdateNFT(network, chain, contract, id, owner string) error
}
