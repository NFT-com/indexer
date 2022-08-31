package failures

import (
	"strings"
)

func TooLarge(err error) bool {
	return strings.Contains(err.Error(), "the message response is too large")
}
