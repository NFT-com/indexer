package failures

import (
	"strings"
)

func TooLarge(err error) bool {
	return err != nil && strings.Contains(err.Error(), "the message response is too large")
}
