package subscriber

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/source"
)

/*type Subscriber interface {
	Subscribe(ctx context.Context, events chan events.Event) error
	Close() error
}*/

type Subscriber struct {
	log zerolog.Logger

	currentSource int
	sources       []source.Source
	parser        block.Parser
}

func NewSubscriber(log zerolog.Logger, parser block.Parser, sources ...source.Source) (*Subscriber, error) {
	s := &Subscriber{
		log:           log.With().Str("component", "subscriber").Logger(),
		currentSource: 0,
		sources:       sources,
		parser:        parser,
	}

	// FIXME: Sanitize input?
	if len(sources) == 0 {
		return nil, errors.New("invalid sources amount")
	}

	return s, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, events chan events.Event) error {
loop:
	for {
		nextBlock := s.sources[s.currentSource].Next()
		if nextBlock == nil {
			s.currentSource++

			if len(s.sources) >= s.currentSource {
				break loop
			}
		}

		blockEvents, err := s.parser.ParseBlock(ctx, nextBlock)
		if err != nil {
			s.log.Error().Str("block", nextBlock.String()).Err(err).Msg("could not parse header")
		}

		for _, event := range blockEvents {
			events <- event
		}
	}

	return nil
}

func (s *Subscriber) Close() error {
	for _, sc := range s.sources {
		err := sc.Close()
		if err != nil {
			s.log.Error().Err(err).Msg("failed to close source")
		}
	}
	return nil
}
