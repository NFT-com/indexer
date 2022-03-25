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

func (s *Store) UpsertSaleEvents(events []events.Sale) error {
	query := s.sqlBuilder.
		Insert(saleEventTableName).
		Columns(saleEventTableColumns...).
		Suffix(saleTableOnConflictStatement)

	for _, event := range events {
		query = query.Values(event.ID, event.Block, event.TransactionHash, event.MarketplaceID, event.Seller, event.Buyer, event.Price, event.EmittedAt)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not insert sale events: %w", err)
	}

	return nil
}
