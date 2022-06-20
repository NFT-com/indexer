package lambda

import (
	"github.com/NFT-com/indexer/models/metadata"
	"github.com/ethereum/go-ethereum/core/types"
)

type LogsFetcher interface {
	Logs(addresses []string, eventTypes []string, from uint64, to uint64) ([]types.Log, error)
}

type URIFetcher interface {
	ERC721(address string, tokenID string) (string, error)
	ERC1155(address string, tokenID string) (string, error)
}

type MetadataFetcher interface {
	Fetch(uri string) (metadata.Token, error)
}
