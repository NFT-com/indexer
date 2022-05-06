package creator

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type CollectionStore interface {
	Combinations(chainID uint64) ([]*jobs.Combination, error)
}

type ParsingStore interface {
	Pending(chainID uint64) (uint, error)
	Latest(chainID uint64, contractAddress string, eventHash string) (*jobs.Parsing, error)
	Insert(parsings *jobs.Parsing) error
}
