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

	// retrieval of token URI for deleted Fighter NFT
	case strings.Contains(msg, "URI query for nonexistent token"):
		return true

	// retrieval of token URI for ENS NFT
	case strings.Contains(msg, "execution reverted"):
		return true

	// decoding of JSON payload for Crypto Kitties NFT
	case strings.Contains(msg, "cannot unmarshal object into Go struct field Token.attributes of type []metadata.Attribute"):
		return true

	// no token metadata exists for token
	case strings.Contains(msg, "no metadata for token"):
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
