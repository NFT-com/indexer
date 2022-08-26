package results

import (
	"strings"
)

func Deleted(err error) bool {

	if err == nil {
		return false
	}

	msg := err.Error()
	switch {

	// retrieval of URI for deleted NFT should not retry
	case strings.Contains(msg, "URI query for nonexistent token"):
		return true

	default:
		return false
	}
}
