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

	// retrieval of URI for deleted NFT
	case strings.Contains(msg, "URI query for nonexistent token"):
		return true

	// unsupported complex OpenSea edge cases
	case strings.Contains(msg, "invalid topic length"):
		return true
	case strings.Contains(msg, "offers are empty"):
		return true
	case strings.Contains(msg, "multiple offers not supported"):
		return true
	case strings.Contains(msg, "considerations are empty"):
		return true
	case strings.Contains(msg, "multiple considerations not supported"):
		return true

	// too many logs returned from node API
	case strings.Contains(msg, "the message response is too large"):
		return true

	// runtime errors that should never happen...
	case strings.Contains(msg, "index out of range"):
		return true
	case strings.Contains(msg, "slice bounds out of range"):
		return true

	default:
		return false
	}
}
