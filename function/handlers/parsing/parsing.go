package parsing

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function/processors/parsing"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/networks/web3"
)

var (
	errParserNotFound = errors.New("parser not found")
)

// Initializer initializes the parser to use with the network client.
type Initializer func(client networks.Network) ([]parsing.Parser, error)

// Handler handles the parsing message from queue.
type Handler struct {
	log         zerolog.Logger
	initializer Initializer
}

// NewHandler creates a new parsing handler consumer.
func NewHandler(log zerolog.Logger, initializer Initializer) *Handler {
	h := Handler{
		log:         log.With().Str("component", "parsing_handler").Logger(),
		initializer: initializer,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, job jobs.Parsing) (interface{}, error) {
	h.log.Debug().
		Str("block", job.BlockNumber).
		Str("event", job.EventType).
		Str("contract", job.Address).
		Msg("processing job")

	network, err := web3.New(ctx, job.ChainURL)
	if err != nil {
		return nil, fmt.Errorf("could not create web3 client: %w", err)
	}
	defer network.Close()

	parser, err := h.getParser(network, job.StandardType)
	if err != nil {
		return nil, fmt.Errorf("could not get parser: %w", err)
	}

	rawLogs, err := network.BlockEvents(ctx, job.BlockNumber, job.EventType, job.Address)
	if err != nil {
		return nil, fmt.Errorf("could not get block events: %w", err)
	}

	logs := make([]log.Log, 0, len(rawLogs))
	for _, rawLog := range rawLogs {
		log, err := parser.ParseRawLog(rawLog)
		if err != nil {
			return nil, fmt.Errorf("could not parse raw event: %w", err)
		}

		logs = append(logs, *log)
	}

	return logs, nil
}

func (h *Handler) getParser(network networks.Network, standardType string) (parsing.Parser, error) {
	parsers, err := h.initializer(network)
	if err != nil {
		return nil, fmt.Errorf("could not initialize parsers: %w", err)
	}

	for _, parser := range parsers {
		if parser.Type() == standardType {
			return parser, nil
		}
	}

	return nil, errParserNotFound
}
