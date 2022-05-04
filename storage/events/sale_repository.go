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

	query := s.build.
		Insert(TableSaleEvents).
		Columns(ColumnsSaleEvents...).
		Suffix(ConflictSaleEvents)

	for _, sale := range sales {
		query = query.Values(
			sale.ID,
			sale.BlockNumber,
			sale.EventIndex,
			sale.TransactionHash,
			sale.MarketplaceAddress,
			sale.SellerAddress,
			sale.BuyerAddress,
			sale.TradePrice,
			sale.EmittedAt,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert sale event: %w",
			err)
	}

	return nil
}
