package parsing

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NFT-com/indexer/function/handlers/parsing"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/models/events"
)

func (d *Parsing) processLogs(input parsing.Input, logs []log.Log) error {
	for _, l := range logs {
		chain, err := d.dataStore.Chain(l.ChainID)
		if err != nil {
			return fmt.Errorf("could not get chain: %w", err)
		}

		if l.NeedsActionJob {
			err := d.jobStore.CreateActionJob(&jobs.Action{
				ID:          uuid.New().String(),
				ChainURL:    input.ChainURL,
				ChainID:     input.ChainID,
				ChainType:   input.ChainType,
				BlockNumber: l.Block,
				Address:     l.Contract,
				Standard:    l.Standard,
				Event:       l.Event,
				TokenID:     l.NftID,
				Type:        l.ActionType.String(),
				Status:      jobs.StatusCreated,
			})
			if err != nil {
				return fmt.Errorf("could not create action job: %w", err)
			}
		}

		switch l.Type {
		case log.Mint:
			collection, err := d.dataStore.Collection(chain.ID, l.Contract, l.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection (chainID:%s contract:%s): %w", chain.ChainID, l.Contract, err)
			}

			event := events.Mint{
				ID:              l.ID,
				CollectionID:    collection.ID,
				Block:           l.Block,
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				TokenID:         l.NftID,
				Owner:           l.ToAddress,
				EmittedAt:       l.EmittedAt,
			}

			err = d.eventStore.UpsertMintEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert mint event: %w", err)
			}

		case log.Sale:
			marketplace, err := d.dataStore.Marketplace(chain.ID, l.Contract)
			if err != nil {
				return fmt.Errorf("could not get marketplace: %w", err)
			}

			event := events.Sale{
				ID:              l.ID,
				MarketplaceID:   marketplace.ID,
				Block:           l.Block,
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				Seller:          l.ToAddress,
				Buyer:           l.FromAddress,
				Price:           l.Price,
				EmittedAt:       l.EmittedAt,
			}

			err = d.eventStore.UpsertSaleEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert sale event: %w", err)
			}

		case log.Transfer:
			collection, err := d.dataStore.Collection(chain.ID, l.Contract, l.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection (chainID:%s contract:%s): %w", chain.ChainID, l.Contract, err)
			}

			event := events.Transfer{
				ID:              l.ID,
				CollectionID:    collection.ID,
				Block:           l.Block,
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				TokenID:         l.NftID,
				FromAddress:     l.FromAddress,
				ToAddress:       l.ToAddress,
				EmittedAt:       l.EmittedAt,
			}

			err = d.eventStore.UpsertTransferEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert transfer event: %w", err)
			}

		case log.Burn:
			collection, err := d.dataStore.Collection(chain.ID, l.Contract, l.ContractCollectionID)
			if err != nil {
				return fmt.Errorf("could not get collection (chainID:%s contract:%s): %w", chain.ChainID, l.Contract, err)
			}

			event := events.Burn{
				ID:              l.ID,
				CollectionID:    collection.ID,
				Block:           l.Block,
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				TokenID:         l.NftID,
				EmittedAt:       l.EmittedAt,
			}

			err = d.eventStore.UpsertBurnEvent(event)
			if err != nil {
				return fmt.Errorf("could not upsert burn event: %w", err)
			}

		default:
			d.log.Error().Str("type", l.Type.String()).Msg("got unexpected log type")
		}
	}

	return nil
}
