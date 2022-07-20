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
	case strings.Contains(msg, "tokenURI: URI query for nonexistent token"):
		return true

	default:
		return false
	}
}
