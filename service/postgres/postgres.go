package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

const defaultParsingJobLimit = 100

// Store represents the database storage struct.
type Store struct {
	// Maximum amount of parsing jobs to keep in DB for
	// each event type and address combination, between
	// each cleanup.
	parsingJobLimit int
	build           squirrel.StatementBuilderType
}

// NewStore returns a new store or error. Requires a database connection.
func NewStore(db *sql.DB, options ...func(*Store)) (*Store, error) {

	s := Store{
		parsingJobLimit: defaultParsingJobLimit,
	}
	for _, opt := range options {
		opt(&s)
	}

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	s.build = squirrel.
		StatementBuilder.
		RunWith(squirrel.NewStmtCache(db)).
		PlaceholderFormat(squirrel.Dollar)

	return &s, nil
}

func WithParsingJobLimit(limit int) func(*Store) {
	return func(store *Store) {
		store.parsingJobLimit = limit
	}
}
