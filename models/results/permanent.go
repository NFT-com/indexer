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

	// retrieval of URI for deleted NFT should not retry
	case strings.Contains(msg, "URI query for nonexistent token"):
		return true
	case strings.Contains(msg, "execution reverted"):
		return true

	// too many logs returned from node API should not retry
	case strings.Contains(msg, "the message response is too large"):
		return true

	// error when parsing event that should be fixed
	case strings.Contains(msg, "invalid number of topics"):
		return true

	// runtime errors are bugs and should always hard fail
	case strings.Contains(msg, "runtime error"):
		return true

	default:
		return false
	}
}
