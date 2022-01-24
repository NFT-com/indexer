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
		strings.ToLower("0x06012c8cf97bead5deae237070f9587f8e7a266d"): cryptokittiesABI,
	}[strings.ToLower(address)]
	if !ok {
		return "", store.ErrNotFound
	}

	return value, nil
}
