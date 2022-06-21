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

	// failures to retrieve data from Ethereum Node
	case strings.Contains(msg, "Too Many Requests"):
		return true
	case strings.Contains(msg, "Internal Server Error"):
		return true
	case strings.Contains(msg, "Bad Gateway"):
		return true
	case strings.Contains(msg, "Gateway Timeout"):
		return true

	// failure due to Lambda file descriptor limit
	case strings.Contains(msg, "too many open files"):
		return true
	case strings.Contains(msg, "no such host"):
		return true

	// failure due to HTTP request issue
	case strings.Contains(msg, "Client.Timeout exceeded while awaiting headers"):
		return true
	case strings.Contains(msg, "bad response code"):
		return true

	// failure due to conflicting SQL transactions
	case strings.Contains(msg, "deadlock detected"):
		return true

	default:
		return false
	}
}
