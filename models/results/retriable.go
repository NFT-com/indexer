package results

import (
	"strings"
)

func Permanent(err error) bool {

	if err == nil {
		return false
	}

	msg := err.Error()
	switch {

	// failure to retrieve token uri
	case strings.Contains(msg, "tokenURI: URI query for nonexistent token"):
		return false

	default:
		return true
	}
}
