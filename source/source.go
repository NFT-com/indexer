package source

import "github.com/NFT-com/indexer/block"

type Source interface {
	Next() *block.Block
	Close() error
}
