package events

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/events"
)

type MintRepository struct {
	build squirrel.StatementBuilderType
}

func NewMintRepository(db *sql.DB) *MintRepository {

	cache := squirrel.NewStmtCache(db)
	m := MintRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &m
}

func (m *MintRepository) Upsert(mint *events.Mint) error {

	_, err := m.build.
		Insert(TableMintEvents).
		Columns(ColumnsMintEvents...).
		Values(
			mint.ID,
			mint.Block,
			mint.EventIndex,
			mint.TransactionHash,
			mint.CollectionID,
			mint.TokenID,
			mint.Owner,
			mint.EmittedAt,
		).
		Suffix(ConflictMintEvents).
		Exec()
	if err != nil {
		return fmt.Errorf("could not upsert mint event: %w",
			err)
	}

	return nil
}
