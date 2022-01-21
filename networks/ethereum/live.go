package ethereum

import (
	"context"

	"github.com/NFT-com/indexer/block"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type LiveSource struct {
	log     zerolog.Logger
	headers chan *types.Header
	done    chan error
}

func NewLive(ctx context.Context, log zerolog.Logger, client *ethclient.Client) (*LiveSource, error) {
	l := LiveSource{
		log:     log.With().Str("component", "live_source").Logger(),
		headers: make(chan *types.Header),
		done:    make(chan error),
	}

	sub, err := client.SubscribeNewHead(ctx, l.headers)
	if err != nil {
		return nil, err
	}

	go func() {
		err = <-sub.Err()
		l.done <- err
	}()

	return &l, nil
}

func (s *LiveSource) Next(ctx context.Context) *block.Block {
	select {
	case header := <-s.headers:
		b := block.Block(header.Hash().Hex())
		return &b

	case err := <-s.done:
		if err != nil {
			s.log.Error().Err(err).Msg("could not subscribe to header")
			break
		}
		s.log.Info().Msg("gracefully stopped")
	case <-ctx.Done():
		s.log.Info().Msg("context closed stopped")
	}

	return nil
}

func (s *LiveSource) Close() error {
	close(s.done)
	return nil
}
