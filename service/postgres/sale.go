package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/events"
)

func (s *Store) UpsertSaleEvent(event events.Sale) error {
	_, err := s.sqlBuilder.
		Insert(saleEventTableName).
		Columns(saleEventTableColumns...).
		Values(event.ID, event.Block, event.TransactionHash, event.MarketplaceID, event.Seller, event.Buyer, event.Price, event.EmittedAt).
		Suffix(saleTableOnConflictStatement).
		Exec()
	if err != nil {
		return fmt.Errorf("could not insert sale event: %w", err)
	}

	return nil
}
