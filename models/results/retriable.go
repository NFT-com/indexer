package results

import (
	"strings"
)

func Retriable(err error) bool {

	if err == nil {
		return false
	}

	msg := err.Error()
	switch {

	case strings.Contains(msg, "Too Many Requests"):
		return true

	case strings.Contains(msg, "too many open files"):
		return true

	case strings.Contains(msg, "no such host"):
		return true

	default:
		return false
	}
}
