package jobs

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/NFT-com/indexer/models/jobs"
)

type FailuresRepository struct {
	build squirrel.StatementBuilderType
}

func NewFailuresRepository(db *sql.DB) *FailuresRepository {

	cache := squirrel.NewStmtCache(db)
	p := FailuresRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &p
}

func (p *FailuresRepository) Parsing(parsing *jobs.Parsing) error {

	_, err := p.build.
		Insert("parsings").
		Columns("id", "chain_id", "contract_addresses", "event_hashes", "start_height", "end_height", "job_status", "input_data").
		Values(
			parsing.ID,
			parsing.ChainID,
			pq.Array(parsing.ContractAddresses),
			pq.Array(parsing.EventHashes),
			parsing.StartHeight,
			parsing.EndHeight,
			parsing.JobStatus,
			parsing.InputData,
		).Exec()
	if err != nil {
		return fmt.Errorf("could not insert parsing job: %w", err)
	}

	return nil
}

func (p *FailuresRepository) Addition(addition *jobs.Addition) (uint, error) {

	result, err := p.build.
		Select("COUNT(id)").
		From("parsings").
		Where("chain_id = ?", chainID).
		Where("job_status != ?", "finished").
		Where("job_status != ?", "failed").
		Query()
	if err != nil {
		return 0, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return 0, fmt.Errorf("could not get result: %w", err)
	}

	if !result.Next() {
		return 0, sql.ErrNoRows
	}

	var count uint
	err = result.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not scan row: %w", err)
	}

	return count, nil
}

func (p *FailuresRepository) Update(chainID uint64, contractAddress string, eventHash string) (uint64, error) {

	result, err := p.build.
		Select("MAX(end_height)").
		From("parsings").
		Where("chain_id = ?", chainID).
		Where("? = ANY(contract_addresses)", contractAddress).
		Where("? = ANY(event_hashes)", eventHash).
		Query()
	if err != nil {
		return 0, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return 0, fmt.Errorf("could not get result: %w", err)
	}

	if !result.Next() {
		return 0, sql.ErrNoRows
	}

	var height sql.NullInt64
	err = result.Scan(&height)
	if err != nil {
		return 0, fmt.Errorf("could not scan row: %w", err)
	}

	if !height.Valid {
		return 0, sql.ErrNoRows
	}

	return uint64(height.Int64), nil
}
