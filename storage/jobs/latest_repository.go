package jobs

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/jobs"
)

type CreatedRepository struct {
	build squirrel.StatementBuilderType
}

func NewCreatedRepository(db *sql.DB) *CreatedRepository {

	cache := squirrel.NewStmtCache(db)
	a := CreatedRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &a
}

func (a *CreatedRepository) Upsert(createds []jobs.Created) error {

	if len(createds) == 0 {
		return nil
	}

	query := a.build.
		Insert("created").
		Columns(
			"chain_id",
			"contract_address",
			"event_hash",
			"block_height",
			"last_id",
			"created_at",
		).
		Suffix("ON CONFLICT(chain_id, contract_address, event_hash) DO UPDATE SET" +
			"block_height = EXCLUDED.block_height, " +
			"last_id = EXCLUDED.last_id, " +
			"created_at = EXCLUDED.created_at")

	for _, created := range createds {
		query = query.Values(
			created.ChainID,
			created.ContractAddress,
			created.EventHash,
			created.BlockHeight,
			created.LastID,
			created.CreatedAt,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not update latest created: %w", err)
	}

	return nil
}

func (a *CreatedRepository) Retrieve() ([]jobs.Created, error) {

	result, err := a.build.
		Select(
			"chain_id",
			"contract_address",
			"event_hash",
			"block_height",
			"last_id",
			"created_at",
		).
		From("created").
		Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve action job: %w", err)
	}
	if !result.Next() {
		return nil, sql.ErrNoRows
	}

	var createds []jobs.Created
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", result.Err())
		}

		var created jobs.Created
		err = result.Scan(
			&created.ChainID,
			&created.ContractAddress,
			&created.EventHash,
			&created.BlockHeight,
			&created.LastID,
			&created.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		createds = append(createds, created)
	}

	return createds, nil
}
