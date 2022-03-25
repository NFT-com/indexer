package consumer

import (
	"fmt"

	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/log"
)

func (d *Parsing) processLogs(logs []log.Log) error {
	for _, l := range logs {
		chain, err := d.store.Chain(l.ChainID)
		if err != nil {
			return fmt.Errorf("could not get chain: %w", err)
		}

		switch l.Type {
		case log.Mint:
			collection, err := d.store.Collection(chain.ID, l.Contract, l.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection: %w", err)
			}

			event := events.Mint{
				ID:              l.ID,
				CollectionID:    collection.ID,
				Block:           l.Block,
				TransactionHash: l.TransactionHash,
				Owner:           l.ToAddress,
			}

			err = d.store.UpsertMintEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert mint event: %w", err)
			}
		case log.Sale:
			marketplace, err := d.store.Marketplace(chain.ID, l.Contract)
			if err != nil {
				return fmt.Errorf("could not get marketplace: %w", err)
			}

			event := events.Sale{
				ID:              l.ID,
				MarketplaceID:   marketplace.ID,
				Block:           l.Block,
				TransactionHash: l.TransactionHash,
				Seller:          l.ToAddress,
				Buyer:           l.FromAddress,
				Price:           l.Price,
			}

			err = d.store.UpsertSaleEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert sale event: %w", err)
			}
		case log.Transfer:
			collection, err := d.store.Collection(chain.ID, l.Contract, l.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection: %w", err)
			}

			event := events.Transfer{
				ID:              l.ID,
				CollectionID:    collection.ID,
				Block:           l.Block,
				TransactionHash: l.TransactionHash,
				FromAddress:     l.FromAddress,
				ToAddress:       l.ToAddress,
			}

			err = d.store.UpsertTransferEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert transfer event: %w", err)
			}
		case log.Burn:
			collection, err := d.store.Collection(chain.ID, l.Contract, l.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection: %w", err)
			}

			event := events.Burn{
				ID:              l.ID,
				CollectionID:    collection.ID,
				Block:           l.Block,
				TransactionHash: l.TransactionHash,
			}

			err = d.store.UpsertBurnEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert burn event: %w", err)
			}
		default:
			d.log.Error().Str("type", l.Type.String()).Msg("got unexpected log type")
		}
	}

	return nil
}
