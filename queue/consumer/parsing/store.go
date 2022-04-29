package parsing

import (
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/models/events"
)

type Store interface {
	UpsertMintEvent(event events.Mint) error
	UpsertSaleEvent(event events.Sale) error
	UpsertTransferEvent(event events.Transfer) error
	UpsertBurnEvent(event events.Burn) error

	Chain(chainID string) (*chain.Chain, error)
	Collection(chainID, address, contractCollectionID string) (*chain.Collection, error)
	Marketplace(chainID, address string) (*chain.Marketplace, error)

	CreateAdditionJob(job *jobs.Addition) error
	ParsingJob(id string) (*jobs.Parsing, error)
	UpdateParsingJobsStatus(ids []string, status jobs.Status) error
}
