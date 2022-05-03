package lambdas

import (
	"github.com/NFT-com/indexer/models/metadata"
)

type MetadataFetcher interface {
	Fetch(uri string) (metadata.Token, error)
}
