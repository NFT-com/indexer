package subscriber

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/source"
)

type Subscriber struct {
	log zerolog.Logger

	currentSource int
	sources       []source.Source
	parser        block.Parser
	done          chan struct{}
}

func NewSubscriber(log zerolog.Logger, parser block.Parser, sources []source.Source) (*Subscriber, error) {
	// FIXME: Sanitize input?
	if len(sources) == 0 {
		return nil, errors.New("invalid sources amount")
	}

	s := Subscriber{
		log:           log.With().Str("component", "subscriber").Logger(),
		currentSource: 0,
		sources:       sources,
		parser:        parser,
		done:          make(chan struct{}),
	}

	return &s, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, events chan *event.Event) error {
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

				blockEvents, err := s.parser.Parse(ctx, nextBlock)
				if err != nil {
					s.log.Error().Str("block", nextBlock.String()).Err(err).Msg("could not parse header")
					continue
				}

				for _, event := range blockEvents {
					events <- event
				}
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
