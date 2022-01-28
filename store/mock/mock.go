package mock

import (
	"context"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
	"github.com/NFT-com/indexer/store"
	"strings"
)

type Mock struct {
	nfts   map[string]*nft.NFT
	events map[string]*event.ParsedEVent
}

func New() *Mock {
	m := Mock{
		nfts:   map[string]*nft.NFT{},
		events: map[string]*event.ParsedEVent{},
	}

	return &m
}

func (m *Mock) SaveNFT(_ context.Context, nft *nft.NFT) error {
	m.nfts[nft.ID] = nft
	return nil
}

func (m *Mock) UpdateNFTOwner(_ context.Context, _, _, _, id, newOwner string) error {
	if _, ok := m.nfts[id]; !ok {
		return store.ErrNotFound
	}

	m.nfts[id].Owner = newOwner
	return nil
}

func (m *Mock) BurnNFT(ctx context.Context, _, _, _, id string) error {
	delete(m.nfts, id)
	return nil
}

func (m *Mock) SaveEvent(_ context.Context, event *event.ParsedEVent) error {
	m.events[event.ID] = event
	return nil
}

func (m *Mock) GetContractType(_ context.Context, _, _, address string) (string, error) {
	value, ok := map[string]string{
		strings.ToLower("0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"): "erc721",
		strings.ToLower("0x06012c8cf97bead5deae237070f9587f8e7a266d"): "custom",
	}[strings.ToLower(address)]
	if !ok {
		return "", store.ErrNotFound
	}

	return value, nil
}

func (m *Mock) GetContractABI(_ context.Context, _, _, address string) (string, error) {
	value, ok := map[string]string{
		strings.ToLower("0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"): erc721ABI,
		strings.ToLower("0x86b18D285C1990Ea16f67D3F22D79970D418C3CE"): erc721ABI,
		strings.ToLower("0x06012c8cf97bead5deae237070f9587f8e7a266d"): cryptokittiesABI,
	}[strings.ToLower(address)]
	if !ok {
		return "", store.ErrNotFound
	}

	return value, nil
}

func (m *Mock) UpdateContractURI(_ context.Context, _, _, _, _ string) error {
	return nil
}
