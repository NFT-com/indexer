package parsing

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/parsers/erc721/transfer"
)

type Handler struct {
	log zerolog.Logger
}

func NewHandler(log zerolog.Logger) *Handler {
	h := Handler{
		log: log.With().Str("component", "parsing_handler").Logger(),
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
		return nil, fmt.Errorf("could create web3 client: %w", err)
	}
	defer network.Close()

	parser := transfer.NewParser()

	rawEvents, err := network.BlockEvents(ctx, job.BlockNumber, job.EventType, job.Address)
	if err != nil {
		return nil, fmt.Errorf("could not get block events: %w", err)
	}

	parsedEvents := make([]events.Event, 0)
	for _, rawEvent := range rawEvents {
		parsedEvent, err := parser.ParseRawEvent(rawEvent)
		if err != nil {
			log.Error().Err(err).Msg("could not parse raw events")
			return nil, fmt.Errorf("could not parse raw events: %w", err)
		}

		parsedEvents = append(parsedEvents, parsedEvent)
	}

	jobResult := jobs.ParsingResult{
		RawEvents:    rawEvents,
		ParsedEvents: parsedEvents,
	}

	return jobResult, nil
}
