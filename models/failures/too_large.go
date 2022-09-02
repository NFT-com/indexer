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

	case strings.Contains(msg, "unexpected EOF"):
		return true

	case strings.Contains(msg, "body size is too long"):
		return true
	}

	return false
}
