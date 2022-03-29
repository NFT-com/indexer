package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

// Store represents the database storage struct.
type Store struct {
	sqlBuilder squirrel.StatementBuilderType
}

// NewStore returns a new store or error. Requires a database connection.
func NewStore(db *sql.DB) (*Store, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	dbCache := squirrel.NewStmtCache(db)

	sqlBuilder := squirrel.StatementBuilder.RunWith(dbCache)
	sqlBuilder = sqlBuilder.PlaceholderFormat(squirrel.Dollar)

	s := Store{
		sqlBuilder: sqlBuilder,
	}

	return &s, nil
}
