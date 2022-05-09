package pipeline

import (
	"github.com/NFT-com/indexer/models/events"
)

type TransferStore interface {
	Upsert(transfers ...*events.Transfer) error
}

type SaleStore interface {
	Upsert(sales ...*events.Sale) error
}
