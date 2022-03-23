package parsing

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/parsers"
)

type Initializer func(client networks.Network) (parsers.Parser, error)

type Handler struct {
	log         zerolog.Logger
	initializer Initializer
}

func NewHandler(log zerolog.Logger, initializer Initializer) *Handler {
	h := Handler{
		log:         log.With().Str("component", "parsing_handler").Logger(),
		initializer: initializer,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, job jobs.Parsing) (interface{}, error) {
	log := h.log.With().
		Str("block", job.BlockNumber).
		Str("event", job.EventType).
		Str("contract", job.Address).
		Logger()

	network, err := web3.New(ctx, job.ChainURL)
	if err != nil {
		return nil, fmt.Errorf("could not create web3 client: %w", err)
	}
	defer network.Close()

	parser, err := h.initializer(network)
	if err != nil {
		return nil, err
	}

	rawEvents, err := network.BlockEvents(ctx, job.BlockNumber, job.EventType, job.Address)
	if err != nil {
		return nil, fmt.Errorf("could not get block event: %w", err)
	}

	parsedEvents := make([]event.Event, 0, len(rawEvents))
	for _, rawEvent := range rawEvents {
		parsedEvent, err := parser.ParseRawEvent(rawEvent)
		if err != nil {
			log.Error().Err(err).Msg("could not parse raw event")
			return nil, fmt.Errorf("could not parse raw event: %w", err)
		}

		parsedEvents = append(parsedEvents, *parsedEvent)
	}

	return parsedEvents, nil
}
