package jobs

import (
	"github.com/Masterminds/squirrel"
)

type UpdateSetter func(query squirrel.UpdateBuilder) squirrel.UpdateBuilder

func SetStatus(status string) UpdateSetter {
	return func(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
		return query.Set("job_status", status)
	}
}

func SetMessage(message string) UpdateSetter {
	return func(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
		return query.Set("status_message", message)
	}
}
