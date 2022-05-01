package jobs

const (
	StatusCreated    = "created"
	StatusQueued     = "queued"
	StatusProcessing = "processing"
	StatusFailed     = "failed"
	StatusFinished   = "finished"
)

func StatusValid(status string) bool {
	switch status {
	case StatusCreated:
		return true
	case StatusQueued:
		return true
	case StatusProcessing:
		return true
	case StatusFailed:
		return true
	case StatusFinished:
		return true
	default:
		return false
	}
}
