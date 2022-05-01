package jobs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/storage/filters"
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

func (p *ParsingRepository) Insert(job *jobs.Parsing) error {

	_, err := p.build.
		Insert(TableParsingJobs).
		Columns(ColumnsParsingJobs...).
		Values(
			job.ID,
			job.ChainURL,
			job.ChainID,
			job.ChainType,
			job.BlockNumber,
			job.Address,
			job.Standard,
			job.Event,
			job.Status,
		).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create parsing job: %w", err)
	}

	return nil
}

func (p *ParsingRepository) Retrieve(parsingID string) (*jobs.Parsing, error) {

	result, err := p.build.
		Select(ColumnsParsingJobs...).
		From(TableParsingJobs).
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

	var job jobs.Parsing
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.Standard,
		&job.Event,
		&job.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parsing job: %w", err)
	}

	return &job, nil
}

func (p *ParsingRepository) Update(parsing *jobs.Parsing) error {

	// TODO: update all columns once we have a better SQL library
	_, err := p.build.
		Update(TableParsingJobs).
		Where("id = ?", parsing.ID).
		Set("status", parsing.Status).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	return nil
}

func (p *ParsingRepository) Find(wheres ...filters.Where) ([]*jobs.Parsing, error) {

	query := p.build.
		Select(ColumnsParsingJobs...).
		From(TableParsingJobs).
		OrderBy("block_number ASC")

	for _, where := range wheres {
		query = query.Where(where())
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	parsings := make([]*jobs.Parsing, 0)
	for result.Next() {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not load next row: %w", err)
		}

		var parsing jobs.Parsing
		err = result.Scan(
			&parsing.ID,
			&parsing.ChainURL,
			&parsing.ChainID,
			&parsing.ChainType,
			&parsing.BlockNumber,
			&parsing.Address,
			&parsing.Standard,
			&parsing.Event,
			&parsing.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		parsings = append(parsings, &parsing)
	}

	return parsings, nil
}
