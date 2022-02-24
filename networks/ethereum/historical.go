package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
)

type HistoricalSource struct {
	log zerolog.Logger

	client    *ethclient.Client
	nextIndex *big.Int
	endIndex  *big.Int
}

func NewHistorical(ctx context.Context, log zerolog.Logger, client *ethclient.Client, startIndex, endHeight int64) (*HistoricalSource, error) {
	h := HistoricalSource{
		log:       log.With().Str("component", "historical_source").Logger(),
		client:    client,
		nextIndex: big.NewInt(startIndex),
		endIndex:  big.NewInt(endHeight),
	}

	latestHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("could not get latest block header")
		return nil, err
	}

	if latestHeader.Number.CmpAbs(h.endIndex) == -1 {
		h.endIndex = latestHeader.Number
	}

	return &h, nil
}

func (s *HistoricalSource) Next(ctx context.Context) *block.Block {
	if s.nextIndex.CmpAbs(s.endIndex) == 0 {
		return nil
	}

	header, err := s.client.HeaderByNumber(ctx, s.nextIndex)
	if err != nil {
		s.log.Error().Err(err).Str("number", s.nextIndex.String()).Msg("could not get block header")
		return nil
	}

	s.nextIndex = s.nextIndex.Add(s.nextIndex, big.NewInt(1))
	b := block.Block{
		Hash: header.Hash().Hex(),
	}

	return &b
}

func (s *HistoricalSource) Close() error {
	return nil
}
