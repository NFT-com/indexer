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

func (m *MintRepository) Upsert(mints ...*events.Mint) error {

	query := m.build.
		Insert(TableMintEvents).
		Columns(ColumnsMintEvents...).
		Suffix(ConflictMintEvents)

	for _, mint := range mints {
		query = query.Values(
			mint.ID,
			mint.BlockNumber,
			mint.EventIndex,
			mint.CollectionAddress,
			mint.TransactionHash,
			mint.TokenID,
			mint.EmittedAt,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not upsert mint event: %w",
			err)
	}

	return nil
}
