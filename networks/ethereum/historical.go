package ethereum

import (
	"context"
	"errors"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
)

type HistoricalSource struct {
	log zerolog.Logger

	client    Client
	nextIndex int64
	endIndex  int64
}

func NewHistorical(ctx context.Context, log zerolog.Logger, client Client, startIndex, endIndex int64) (*HistoricalSource, error) {
	if client == nil {
		return nil, errors.New("invalid ethereum client")
	}

	h := HistoricalSource{
		log:       log.With().Str("component", "historical_source").Logger(),
		client:    client,
		nextIndex: startIndex,
		endIndex:  endIndex,
	}

	latestHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		h.log.Error().Err(err).Msg("could not get latest block header")
		return nil, err
	}

	if endIndex > latestHeader.Number.Int64() {
		h.endIndex = latestHeader.Number.Int64()
	}

	return &h, nil
}

func (s *HistoricalSource) Next(ctx context.Context) *block.Block {
	if s.nextIndex > s.endIndex {
		return nil
	}

	header, err := s.client.HeaderByNumber(ctx, big.NewInt(s.nextIndex))
	if err != nil {
		s.log.Error().Err(err).Int64("header", s.nextIndex).Msg("could not get block header")
		return nil
	}

	s.nextIndex++
	b := block.Block(header.Hash().Hex())
	return &b
}

func (s *HistoricalSource) Close() error {
	return nil
}
