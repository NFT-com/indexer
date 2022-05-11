package jobs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/NFT-com/indexer/models/jobs"
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

func (a *ActionRepository) Insert(actions ...*jobs.Action) error {

	if len(actions) == 0 {
		return nil
	}

	query := a.build.
		Insert("actions").
		Columns("id", "chain_id", "contract_address", "token_id", "action_type", "block_height", "job_status", "input_data")

	for _, action := range actions {
		query = query.Values(
			action.ID,
			action.ChainID,
			action.ContractAddress,
			action.TokenID,
			action.ActionType,
			action.BlockHeight,
			action.JobStatus,
			action.InputData,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not create action job: %w", err)
	}

	return nil
}

func (a *ActionRepository) Retrieve(actionID string) (*jobs.Action, error) {

	result, err := a.build.
		Select("id", "chain_id", "contract_address", "token_id", "action_type", "block_height", "job_status", "input_data").
		From("actions").
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

	var action jobs.Action
	err = result.Scan(
		&action.ID,
		&action.ChainID,
		&action.ContractAddress,
		&action.TokenID,
		&action.ActionType,
		&action.BlockHeight,
		&action.JobStatus,
		&action.InputData,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve action job: %w", err)
	}

	return &action, nil
}

func (a *ActionRepository) UpdateStatus(status string, statusMessage string, actionIDs ...string) error {

	query := a.build.
		Update("actions").
		Where("id = ANY(?)", pq.Array(actionIDs)).
		Set("job_status", status).
		Set("updated_at", time.Now())

	if statusMessage != "" {
		query = query.Set("status_message", statusMessage)
	}

	result, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not update action job status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not count affected rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *ActionRepository) List(status string) ([]*jobs.Action, error) {

	result, err := a.build.
		Select("id", "chain_id", "contract_address", "token_id", "action_type", "block_height", "job_status", "input_data").
		From("actions").
		Where("job_status = ?", status).
		OrderBy("block_height ASC").
		Query()
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
			&action.ChainID,
			&action.ContractAddress,
			&action.TokenID,
			&action.ActionType,
			&action.BlockHeight,
			&action.JobStatus,
			&action.InputData,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan next row: %w", err)
		}

		actions = append(actions, &action)
	}

	return actions, nil
}
