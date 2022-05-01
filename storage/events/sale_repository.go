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

func (s *SaleRepository) Upsert(sale events.Sale) error {

	_, err := s.build.
		Insert(TableSaleEvents).
		Columns(ColumnsSaleEvents...).
		Values(
			sale.ID,
			sale.Block,
			sale.EventIndex,
			sale.TransactionHash,
			sale.MarketplaceID,
			sale.Seller,
			sale.Buyer,
			sale.Price,
			sale.EmittedAt,
		).
		Suffix(ConflictSaleEvents).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert sale event: %w",
			err)
	}

	return nil
}
