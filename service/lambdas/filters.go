package lambdas

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func FilterForTransactionHash(logs []types.Log, transaction string) []types.Log {
	filtered := make([]types.Log, 0, len(logs))

	for _, log := range logs {
		if log.TxHash.Hex() == transaction {
			filtered = append(filtered, log)
		}
	}

	return filtered
}
