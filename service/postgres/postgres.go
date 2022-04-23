package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

// Store represents the database storage struct.
type Store struct {
	build squirrel.StatementBuilderType
}

// NewStore returns a new store or error. Requires a database connection.
func NewStore(db *sql.DB) (*Store, error) {

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	build := squirrel.
		StatementBuilder.
		RunWith(squirrel.NewStmtCache(db)).
		PlaceholderFormat(squirrel.Dollar)

	s := Store{
		build: build,
	}

	return &s, nil
}
