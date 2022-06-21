package jobs

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

type BoundaryRepository struct {
	build squirrel.StatementBuilderType
}

func NewBoundaryRepository(db *sql.DB) *BoundaryRepository {

	cache := squirrel.NewStmtCache(db)
	b := BoundaryRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &b
}

func (b *BoundaryRepository) ForCombination(chainID uint64, address string, event string) (uint64, error) {

	result, err := b.build.
		Select("last_height").
		From("boundaries").
		Where("chain_id = ?", chainID).
		Where("LOWER(contract_address) = LOWER(?)", address).
		Where("event_hash = ?", event).
		Query()
	if err != nil {
		return 0, fmt.Errorf("could not query collection: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return 0, fmt.Errorf("could not get row: %w", result.Err())
	}
	if !result.Next() {
		return 0, sql.ErrNoRows
	}

	var height uint64
	err = result.Scan(&height)
	if err != nil {
		return 0, fmt.Errorf("could not scan row: %w", err)
	}

	return height, nil

}

func (b *BoundaryRepository) Upsert(chainID uint64, addresses []string, events []string, height uint64, jobID string) error {

	if len(addresses) == 0 || len(events) == 0 {
		return nil
	}

	query := b.build.
		Insert("boundaries").
		Columns(
			"chain_id",
			"contract_address",
			"event_hash",
			"last_height",
			"last_id",
			"created_at",
		).
		Suffix("ON CONFLICT (chain_id, contract_address, event_hash) DO UPDATE SET " +
			"last_height = EXCLUDED.last_height, " +
			"last_id = EXCLUDED.last_id, " +
			"updated_at = NOW()")

	for _, address := range addresses {
		for _, event := range events {
			query = query.Values(
				chainID,
				address,
				event,
				height,
				jobID,
				"NOW()",
			)
		}
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert boundaries: %w", err)
	}

	return nil
}
