package mock

import (
	"context"
	"strings"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
	"github.com/NFT-com/indexer/store"
)

type Mock struct {
	log zerolog.Logger
}

func New(log zerolog.Logger) *Mock {
	m := Mock{
		log: log,
	}

	return &m
}

func (m *Mock) SaveNFT(_ context.Context, nft *nft.NFT) error {
	m.log.Info().Interface("nft", *nft).Msg("new nft inserted")
	return nil
}

func (m *Mock) UpdateNFTOwner(_ context.Context, _, _, _, id, newOwner string) error {
	m.log.Info().Str("id", id).Str("newOwner", newOwner).Msg("nft owner updated")
	return nil
}

func (m *Mock) BurnNFT(ctx context.Context, _, _, _, id string) error {
	m.log.Info().Str("id", id).Msg("nft burnt")
	return nil
}

func (m *Mock) SaveEvent(_ context.Context, event *event.ParsedEvent) error {
	m.log.Info().Interface("id", *event).Msg("new event saved")
	return nil
}

func (m *Mock) GetContractType(_ context.Context, _, _, address string) (string, error) {
	value, ok := map[string]string{
		strings.ToLower("0x86b18D285C1990Ea16f67D3F22D79970D418C3CE"): "erc721",
		strings.ToLower("0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"): "custom",
		strings.ToLower("0x06012c8cf97bead5deae237070f9587f8e7a266d"): "custom",
	}[strings.ToLower(address)]
	if !ok {
		return "", store.ErrNotFound
	}

	return value, nil
}

func (m *Mock) GetContractABI(_ context.Context, _, _, address string) (string, error) {
	value, ok := map[string]string{
		strings.ToLower("0x86b18D285C1990Ea16f67D3F22D79970D418C3CE"): erc721ABI,
		strings.ToLower("0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"): ensABI,
		strings.ToLower("0x06012c8cf97bead5deae237070f9587f8e7a266d"): cryptokittiesABI,
	}[strings.ToLower(address)]
	if !ok {
		return "", store.ErrNotFound
	}

	return value, nil
}

func (m *Mock) UpdateContractURI(_ context.Context, _, _, address, uri string) error {
	m.log.Info().Str("contract", address).Str("uri", uri).Msg("contract uri updated")
	return nil
}
