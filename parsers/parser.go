package parsers

import (
	"github.com/NFT-com/indexer/log"
)

type Parser interface {
	ParseRawLog(log log.RawLog) (*log.Log, error)
}
