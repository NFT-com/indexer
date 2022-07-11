package events

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

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

	if len(sales) == 0 {
		return nil
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
			"trade_price",
			"payment_token",
			"emitted_at",
		).
		Suffix("ON CONFLICT DO NOTHING")

	for _, sale := range sales {
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
			sale.TradePrice,
			sale.PaymentToken,
			sale.EmittedAt,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert sale event: %w", err)
	}

	return nil
}

func (s *SaleRepository) Update(sales ...*events.Sale) error {

	for _, sale := range sales {

		query := s.build.
			Update("sales").
			Set("collection_address", sale.CollectionAddress).
			Set("token_id", sale.TokenID).
			Where("id = ?", sale.ID)

		_, err := query.Exec()
		if err != nil {
			return fmt.Errorf("could not update sale event (id: %s): %w", sale.ID, err)
		}

	}

	return nil
}
