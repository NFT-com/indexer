package mocks

import (
	"context"
	"testing"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
)

type Store struct {
	SaveNFTFunc           func(ctx context.Context, nft *nft.NFT) error
	UpdateNFTOwnerFunc    func(ctx context.Context, network, chain, address, id, newOwner string) error
	UpdateNFTMetadataFunc func(ctx context.Context, network, chain, address, id string, data map[string]interface{}) error
	BurnNFTFunc           func(ctx context.Context, network, chain, address, id string) error
	SaveEventFunc         func(ctx context.Context, event *event.ParsedEvent) error
	GetContractTypeFunc   func(ctx context.Context, network, chain, address string) (string, error)
	GetContractABIFunc    func(ctx context.Context, network, chain, address string) (string, error)
	UpdateContractURIFunc func(ctx context.Context, network, chain, address, uri string) error
}

func BaselineStore(t *testing.T) *Store {
	t.Helper()

	c := Store{
		SaveNFTFunc: func(context.Context, *nft.NFT) error {
			return nil
		},
		UpdateNFTOwnerFunc: func(context.Context, string, string, string, string, string) error {
			return nil
		},
		UpdateNFTMetadataFunc: func(context.Context, string, string, string, string, map[string]interface{}) error {
			return nil
		},
		BurnNFTFunc: func(context.Context, string, string, string, string) error {
			return nil
		},
		SaveEventFunc: func(context.Context, *event.ParsedEvent) error {
			return nil
		},
		GetContractTypeFunc: func(context.Context, string, string, string) (string, error) {
			return GenericContractType, nil
		},
		GetContractABIFunc: func(context.Context, string, string, string) (string, error) {
			return GenericContractABI, nil
		},
		UpdateContractURIFunc: func(context.Context, string, string, string, string) error {
			return nil
		},
	}

	return &c
}

func (s *Store) SaveNFT(ctx context.Context, nft *nft.NFT) error {
	return s.SaveNFTFunc(ctx, nft)
}

func (s *Store) UpdateNFTOwner(ctx context.Context, network, chain, address, id, newOwner string) error {
	return s.UpdateNFTOwnerFunc(ctx, network, chain, address, id, newOwner)
}

func (s *Store) UpdateNFTMetadata(ctx context.Context, network, chain, address, id string, data map[string]interface{}) error {
	return s.UpdateNFTMetadataFunc(ctx, network, chain, address, id, data)
}

func (s *Store) BurnNFT(ctx context.Context, network, chain, address, id string) error {
	return s.BurnNFTFunc(ctx, network, chain, address, id)
}

func (s *Store) SaveEvent(ctx context.Context, event *event.ParsedEvent) error {
	return s.SaveEventFunc(ctx, event)
}

func (s *Store) GetContractType(ctx context.Context, network, chain, address string) (string, error) {
	return s.GetContractTypeFunc(ctx, network, chain, address)
}

func (s *Store) GetContractABI(ctx context.Context, network, chain, address string) (string, error) {
	return s.GetContractABIFunc(ctx, network, chain, address)
}

func (s *Store) UpdateContractURI(ctx context.Context, network, chain, address, uri string) error {
	return s.UpdateContractURIFunc(ctx, network, chain, address, uri)
}
