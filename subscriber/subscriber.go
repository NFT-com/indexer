package subscriber

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/source"
)

type Subscriber struct {
	log zerolog.Logger

	currentSource int
	sources       []source.Source
	parser        block.Parser
	done          chan struct{}
}

func NewSubscriber(log zerolog.Logger, sources ...source.Source) (*Subscriber, error) {
	if len(sources) == 0 {
		return nil, errors.New("invalid sources amount")
	}

	s := Subscriber{
		log:           log.With().Str("component", "subscriber").Logger(),
		currentSource: 0,
		sources:       sources,
		done:          make(chan struct{}),
	}

	return &s, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, events chan *block.Block) error {
	for {
		select {
		case <-s.done:
			return nil
		default:
			nextBlock := s.sources[s.currentSource].Next(ctx)
			if nextBlock == nil {
				s.currentSource++
				if s.currentSource >= len(s.sources) {
					return nil
				}

				continue
			}

			events <- nextBlock
		}
	}
}

func (s *Subscriber) Close() error {
	close(s.done)

	for _, sc := range s.sources {
		err := sc.Close()
		if err != nil {
			s.log.Error().Err(err).Msg("could not close source")
		}
	}

	return nil
}
