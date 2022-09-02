package storage

import (
	"github.com/Masterminds/squirrel"
)

type Retrier interface {
	Insert(query squirrel.InsertBuilder) error
	Update(query squirrel.UpdateBuilder) error
	Delete(query squirrel.DeleteBuilder) error
}
