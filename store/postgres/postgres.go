package postgres

import (
	"context"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
)

type Store struct {
}

func (m *Store) SaveNFT(_ context.Context, externalNFT *nft.NFT) error {
	_, err := FromExternalNFT(externalNFT)
	if err != nil {
		return err
	}

	return nil
}

func (m *Store) UpdateNFTOwner(_ context.Context, _, _, _, id, newOwner string) error {
	return nil
}

func (m *Store) BurnNFT(ctx context.Context, _, _, _, id string) error {
	return nil
}

func (m *Store) SaveEvent(_ context.Context, event *event.ParsedEvent) error {
	_, err := FromExternalEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (m *Store) GetContractType(_ context.Context, _, _, address string) (string, error) {
	return "", nil
}

func (m *Store) GetContractABI(_ context.Context, _, _, address string) (string, error) {
	return "", nil
}

func (m *Store) UpdateContractURI(_ context.Context, _, _, _, _ string) error {
	return nil
}
