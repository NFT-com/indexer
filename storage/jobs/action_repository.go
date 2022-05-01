package jobs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/storage/filters"
)

type ActionRepository struct {
	build squirrel.StatementBuilderType
}

func NewActionRepository(db *sql.DB) *ActionRepository {

	cache := squirrel.NewStmtCache(db)
	a := ActionRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &a
}

func (a *ActionRepository) Insert(action *jobs.Action) error {

	_, err := a.build.
		Insert(TableActionJobs).
		Columns(ColumnsActionJobs...).
		Values(
			action.ID,
			action.ChainURL,
			action.ChainID,
			action.ChainType,
			action.BlockNumber,
			action.Address,
			action.Standard,
			action.Event,
			action.TokenID,
			action.ToAddress,
			action.Type,
			action.Status,
		).
		Exec()
	if err != nil {
		return fmt.Errorf("could not create action job: %w", err)
	}

	return nil
}

func (a *ActionRepository) Retrieve(actionID string) (*jobs.Action, error) {

	result, err := a.build.
		Select(ColumnsActionJobs...).
		From(TableActionJobs).
		Where("id = ?", actionID).
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

	var job jobs.Action
	err = result.Scan(
		&job.ID,
		&job.ChainURL,
		&job.ChainID,
		&job.ChainType,
		&job.BlockNumber,
		&job.Address,
		&job.Standard,
		&job.Event,
		&job.TokenID,
		&job.ToAddress,
		&job.Type,
		&job.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve action job: %w", err)
	}

	return &job, nil
}

func (a *ActionRepository) Update(action *jobs.Action) error {

	_, err := a.build.
		Update(TableActionJobs).
		Where("id = ?", action.ID).
		Set("status", action.Status).
		Set("updated_at", time.Now()).
		Exec()
	if err != nil {
		return fmt.Errorf("could not update action job status: %w", err)
	}

	return nil
}

func (a *ActionRepository) Find(wheres ...filters.Where) ([]*jobs.Action, error) {

	query := a.build.
		Select(ColumnsActionJobs...).
		From(TableActionJobs).
		OrderBy("block_number ASC")

	for _, where := range wheres {
		query = query.Where(where())
	}

	result, err := query.Query()
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer result.Close()

	actions := make([]*jobs.Action, 0)
	for result.Next() && result.Err() == nil {

		if result.Err() != nil {
			return nil, fmt.Errorf("could not load next row: %w", err)
		}

		var action jobs.Action
		err = result.Scan(
			&action.ID,
			&action.ChainURL,
			&action.ChainID,
			&action.ChainType,
			&action.BlockNumber,
			&action.Address,
			&action.Standard,
			&action.Event,
			&action.TokenID,
			&action.ToAddress,
			&action.Type,
			&action.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		actions = append(actions, &action)
	}

	return actions, nil
}
