package store

import (
	"context"
	"errors"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storer interface {
	SaveNFT(ctx context.Context, nft *nft.NFT) error
	SaveEvent(ctx context.Context, event *event.ParsedEVent) error

	GetContractType(ctx context.Context, network, chain, address string) (string, error)
	GetContractABI(ctx context.Context, network, chain, address string) (string, error)
}
