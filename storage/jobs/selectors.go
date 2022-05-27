package jobs

import (
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type UpdateSelector func(query squirrel.UpdateBuilder) squirrel.UpdateBuilder

func One(id string) UpdateSelector {
	return func(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
		return query.Where("id = ?", id)
	}
}

func Many(ids []string) UpdateSelector {
	return func(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
		return query.Where("id = ANY(?)", pq.Array(ids))
	}
}
