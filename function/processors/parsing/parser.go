package parsing

import (
	"github.com/NFT-com/indexer/log"
)

type Parser interface {
	Type() string
	ParseRawLog(log log.RawLog) (*log.Log, error)
}
