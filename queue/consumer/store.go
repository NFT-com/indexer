package consumer

import (
	"github.com/NFT-com/indexer/chain"
	"github.com/NFT-com/indexer/collection"
	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/marketplace"
)

type Store interface {
	UpsertMintEvent(event events.Mint) error
	UpsertMintEvents(events []events.Mint) error

	UpsertSaleEvent(event events.Sale) error
	UpsertSaleEvents(events []events.Sale) error

	UpsertTransferEvent(event events.Transfer) error
	UpsertTransferEvents(events []events.Transfer) error

	UpsertBurnEvent(event events.Burn) error
	UpsertBurnEvents(events []events.Burn) error

	Chain(chainID string) (*chain.Chain, error)

	Collection(chainID, address, contractCollectionID string) (*collection.Collection, error)

	Marketplace(chainID, address string) (*marketplace.Marketplace, error)
}
