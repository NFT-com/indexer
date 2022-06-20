package jobs

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/jobs"
)

type FailureRepository struct {
	build squirrel.StatementBuilderType
}

func NewFailureRepository(db *sql.DB) *FailureRepository {

	cache := squirrel.NewStmtCache(db)
	f := FailureRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &f
}

func (b *FailureRepository) Parsing(parsing *jobs.Parsing, message string) error {

	query := b.build.
		Insert("parsing_failures").
		Columns(
			"id",
			"chain_id",
			"start_height",
			"end_height",
			"contract_addresses",
			"event_hashes",
			"failure_message",
		).
		Values(
			parsing.ID,
			parsing.ChainID,
			parsing.StartHeight,
			parsing.EndHeight,
			parsing.ContractAddresses,
			parsing.EventHashes,
			message,
		)

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not insert parsing failure: %w", err)
	}

	return nil
}

func (b *FailureRepository) Addition(addition *jobs.Addition, message string) error {

	query := b.build.
		Insert("parsing_failures").
		Columns(
			"id",
			"chain_id",
			"contract_address",
			"token_id",
			"token_standard",
			"owner_address",
			"token_count",
			"failure_message",
		).
		Values(
			addition.ID,
			addition.ChainID,
			addition.ContractAddress,
			addition.TokenID,
			addition.TokenStandard,
			addition.OwnerAddress,
			addition.TokenCount,
			message,
		)

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not insert addition failure: %w", err)
	}

	return nil
}
