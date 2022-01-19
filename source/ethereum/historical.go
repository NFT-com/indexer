package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/NFT-com/indexer/parse"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type HistoricalSource struct {
	log zerolog.Logger

	client    *ethclient.Client
	nextIndex int64
	endIndex  int64
}

func NewHistorical(ctx context.Context, log zerolog.Logger, client *ethclient.Client, startIndex, endIndex int64) (*HistoricalSource, error) {
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

func (s *HistoricalSource) Next() *parse.Block {
	// FIXME: Should this be here? Or Next should take the ctx?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	header, err := s.client.HeaderByNumber(ctx, big.NewInt(s.nextIndex))
	if err != nil {
		s.log.Error().Err(err).Int64("header", s.nextIndex).Msg("could not get block header")
		return nil
	}

	s.nextIndex++
	b := parse.Block(header.Hash().Hex())
	return &b
}

func (s *HistoricalSource) Close() error {
	return nil
}
