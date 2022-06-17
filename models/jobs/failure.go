package jobs

import (
	"time"
)

type Failure struct {
	JobID      string    `json:"job_id"`
	Type       string    `json:"type"`
	Parameters string    `json:"parameters"`
	CreatedAt  time.Time `json:"created_at"`
	FailedAt   time.Time `json:"failed_at"`
}
