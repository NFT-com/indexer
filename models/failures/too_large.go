package failures

import (
	"strings"
)

func TooLarge(err error) bool {

	msg := err.Error()
	switch {

	case strings.Contains(msg, "the message response is too large"):
		return true

	case strings.Contains(msg, "request timed out"):
		return true
	}

	return false
}
