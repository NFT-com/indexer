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

	// hitting the request limit on Lambda or Ethereum JSON RPC API
	case strings.Contains(msg, "Too Many Requests"):
		return true

	// hitting the Lambda file descriptor limit
	case strings.Contains(msg, "too many open files"):
		return true

	// also hitting the Lambda file descriptor limit;
	// sometimes the DNS queries fail then
	case strings.Contains(msg, "no such host"):
		return true

	// hitting the HTTP timeout when fetching NFT information
	case strings.Contains(msg, "Client.Timeout exceeded while awaiting headers"):
		return true

	// hitting two conflicting transactions on Postgres
	case strings.Contains(msg, "deadlock detected"):
		return true

	// failing to retrieve token details over HTTP
	case strings.Contains(msg, "bad response code"):
		return true

	default:
		return false
	}
}
