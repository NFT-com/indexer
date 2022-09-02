package jobs

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/jobs"
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

func (b *BoundaryRepository) All() ([]*jobs.Boundary, error) {

	result, err := b.build.
		Select(
			"chain_id",
			"contract_address",
			"event_hash",
			"next_height",
		).
		From("boundaries").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	var boundaries []*jobs.Boundary
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var boundary jobs.Boundary
		err = result.Scan(
			&boundary.ChainID,
			&boundary.ContractAddress,
			&boundary.EventHash,
			&boundary.NextHeight,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		boundaries = append(boundaries, &boundary)
	}

	return boundaries, nil
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
			"next_height",
			"last_id",
			"updated_at",
		).
		Suffix("ON CONFLICT (chain_id, contract_address, event_hash) DO UPDATE SET " +
			"next_height = EXCLUDED.next_height, " +
			"last_id = EXCLUDED.last_id, " +
			"updated_at = EXCLUDED.updated_at")

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
