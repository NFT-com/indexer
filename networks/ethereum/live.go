package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
)

type LiveSource struct {
	log     zerolog.Logger
	headers chan *types.Header
	done    chan error
}

func NewLive(ctx context.Context, log zerolog.Logger, client Client) (*LiveSource, error) {
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
		select {
		case err := <-sub.Err():
			l.done <- err
		case <-l.done:
			sub.Unsubscribe()
			close(l.headers)
			return
		}
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
		s.log.Info().Msg("interrupted")
	}

	return nil
}

func (s *LiveSource) Close() error {
	close(s.done)
	return nil
}
