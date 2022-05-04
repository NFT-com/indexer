package pipeline

import (
	"github.com/NFT-com/indexer/models/events"
)

type BurnStore interface {
	Upsert(burns ...*events.Burn) error
}

type MintStore interface {
	Upsert(mints ...*events.Mint) error
}

type TransferStore interface {
	Upsert(transfers ...*events.Transfer) error
}

type SaleStore interface {
	Upsert(sales ...*events.Sale) error
}
