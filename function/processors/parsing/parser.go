package parsing

import (
	"github.com/NFT-com/indexer/log"
)

type Parser interface {
	Type() string
	ParseRawLog(log log.RawLog, standards map[string]string) (*log.Log, error)
}
