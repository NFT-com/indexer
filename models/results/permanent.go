package results

import (
	"strings"
)

func Permanent(err error) bool {

	if err == nil {
		return true
	}

	msg := err.Error()
	switch {

	// failure to retrieve token uri
	case strings.Contains(msg, "URI query for nonexistent token"):
		return true

	// unsupported OpenSea events
	case strings.Contains(msg, "multiple offers not supported"):
		return true
	case strings.Contains(msg, "considerations are empty"):
		return true

	default:
		return false
	}
}
