package jobs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/NFT-com/indexer/models/jobs"
)

type ParsingRepository struct {
	build squirrel.StatementBuilderType
}

func NewParsingRepository(db *sql.DB) *ParsingRepository {

	cache := squirrel.NewStmtCache(db)
	p := ParsingRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &p
}

func (p *ParsingRepository) Insert(parsing *jobs.Parsing) error {

	_, err := p.build.
		Insert("parsings").
		Columns(ColumnsParsingJobs...).
		Values(
			parsing.ID,
			parsing.ChainID,
			parsing.ContractAddresses,
			parsing.EventHashes,
			parsing.StartHeight,
			parsing.EndHeight,
			parsing.Status,
		).Exec()
	if err != nil {
		return fmt.Errorf("could not insert parsing job: %w", err)
	}

	return nil
}

func (p *ParsingRepository) Pending(chainID uint64) (uint, error) {

	result, err := p.build.
		Select("COUNT(id)").
		From("parsings").
		Where("chain_id = ?", chainID).
		Where("status != finished").
		Where("status != failed").
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

func (p *ParsingRepository) Latest(chainID uint64, contractAddress string, eventHash string) (uint64, error) {

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

	var height uint64
	err = result.Scan(&height)
	if err != nil {
		return 0, fmt.Errorf("could not scan row: %w", err)
	}

	return height, nil
}

func (p *ParsingRepository) Retrieve(parsingID string) (*jobs.Parsing, error) {

	result, err := p.build.
		Select(ColumnsParsingJobs...).
		From("parsings").
		Where("id = ?", parsingID).
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}
	if !result.Next() {
		return nil, sql.ErrNoRows
	}

	var parsing jobs.Parsing
	err = result.Scan(
		&parsing.ID,
		&parsing.ChainID,
		&parsing.ContractAddresses,
		&parsing.EventHashes,
		&parsing.StartHeight,
		&parsing.EndHeight,
		&parsing.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}

	return &parsing, nil
}

func (p *ParsingRepository) UpdateStatus(status string, parsingIDs ...string) error {

	_, err := p.build.
		Update("parsings").
		Where("id IN (?)", pq.Array(parsingIDs)).
		Set("status", status).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	return nil
}

func (p *ParsingRepository) Find(wheres ...string) ([]*jobs.Parsing, error) {

	query := p.build.
		Select(ColumnsParsingJobs...).
		From("parsings").
		OrderBy("block_number ASC")

	for _, where := range wheres {
		query = query.Where(where)
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	parsings := make([]*jobs.Parsing, 0)
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not get next row: %w", err)
		}

		var parsing jobs.Parsing
		err = result.Scan(
			&parsing.ID,
			&parsing.ChainID,
			&parsing.ContractAddresses,
			&parsing.EventHashes,
			&parsing.StartHeight,
			&parsing.EndHeight,
			&parsing.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		parsings = append(parsings, &parsing)
	}

	return parsings, nil
}
