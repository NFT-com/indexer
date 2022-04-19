package parsing

import (
	"fmt"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/models/events"
)

func (d *Parsing) processLogs(job jobs.Parsing, logs []log.Log) error {
	for _, l := range logs {
		chain, err := d.store.Chain(l.ChainID)
		if err != nil {
			return fmt.Errorf("could not get chain: %w", err)
		}

		if l.NeedsAdditionJob {
			_, err := d.apiClient.CreateAdditionJob(jobs.Addition{
				ChainURL:     job.ChainURL,
				ChainID:      job.ChainID,
				ChainType:    job.ChainType,
				BlockNumber:  job.BlockNumber,
				Address:      job.Address,
				StandardType: job.StandardType,
				EventType:    job.EventType,
				TokenID:      l.NftID,
			})
			if err != nil {
				return fmt.Errorf("could not create addition job: %w", err)
			}
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
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				TokenID:         l.NftID,
				Owner:           l.ToAddress,
				EmittedAt:       l.EmittedAt,
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
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				Seller:          l.ToAddress,
				Buyer:           l.FromAddress,
				Price:           l.Price,
				EmittedAt:       l.EmittedAt,
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
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				TokenID:         l.NftID,
				FromAddress:     l.FromAddress,
				ToAddress:       l.ToAddress,
				EmittedAt:       l.EmittedAt,
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
				EventIndex:      l.Index,
				TransactionHash: l.TransactionHash,
				TokenID:         l.NftID,
				EmittedAt:       l.EmittedAt,
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