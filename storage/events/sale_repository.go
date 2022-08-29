package events

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/database"
	"github.com/NFT-com/indexer/models/events"
)

type SaleRepository struct {
	build squirrel.StatementBuilderType
}

func NewSaleRepository(db *sql.DB) *SaleRepository {

	cache := squirrel.NewStmtCache(db)
	s := SaleRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &s
}

func (s *SaleRepository) Upsert(sales ...*events.Sale) error {

	for start := 0; start > len(sales); start += database.BatchSize {

		end := start + database.BatchSize
		if end > len(sales) {
			end = len(sales)
		}

		batch := sales[start:end]
		if len(batch) == 0 {
			continue
		}

		query := s.build.
			Insert("sales").
			Columns(
				"id",
				"chain_id",
				"marketplace_address",
				"collection_address",
				"token_id",
				"block_number",
				"transaction_hash",
				"event_index",
				"seller_address",
				"buyer_address",
				"token_count",
				"currency_value",
				"currency_address",
				"emitted_at",
			).
			Suffix("ON CONFLICT DO NOTHING")

		for _, sale := range batch {
			query = query.Values(
				sale.ID,
				sale.ChainID,
				sale.MarketplaceAddress,
				sale.CollectionAddress,
				sale.TokenID,
				sale.BlockNumber,
				sale.TransactionHash,
				sale.EventIndex,
				sale.SellerAddress,
				sale.BuyerAddress,
				sale.TokenCount,
				sale.CurrencyValue,
				sale.CurrencyAddress,
				sale.EmittedAt,
			)
		}

		_, err := query.Exec()
		if err != nil {
			return fmt.Errorf("could not upsert sales batch (start: %d, end: %d): %w", start, end, err)
		}
	}

	return nil
}

func (s *SaleRepository) Update(sales ...*events.Sale) error {

	for _, sale := range sales {

		query := s.build.
			Update("sales").
			Set("collection_address", sale.CollectionAddress).
			Set("token_id", sale.TokenID).
			Set("token_count", sale.TokenCount).
			Set("currency_value", sale.CurrencyValue).
			Set("currency_address", sale.CurrencyAddress).
			Where("id = ?", sale.ID)

		_, err := query.Exec()
		if err != nil {
			return fmt.Errorf("could not update sale event (id: %s): %w", sale.ID, err)
		}

	}

	return nil
}
