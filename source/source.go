package source

import (
	"github.com/NFT-com/indexer/parse"
)

type Source interface {
	Next() *parse.Block
	Close() error
}
