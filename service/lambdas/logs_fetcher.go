package lambdas

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type LogsFetcher interface {
	Logs(addresses []string, eventTypes []string, from uint64, to uint64) ([]types.Log, error)
}
