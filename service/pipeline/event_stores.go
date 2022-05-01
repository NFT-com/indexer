package pipeline

import (
	"github.com/NFT-com/indexer/models/events"
)

type MintStore interface {
	Upsert(mint events.Mint) error
}

type TransferStore interface {
	Upsert(transfer events.Transfer) error
}

type SaleStore interface {
	Upsert(sale events.Sale) error
}

type BurnStore interface {
	Upsert(burn events.Burn) error
}
