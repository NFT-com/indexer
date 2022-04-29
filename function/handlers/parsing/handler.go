package parsing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function/processors/parsing"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/networks/web3"
)

var (
	errParserNotFound = errors.New("parser not found")
)

type Input struct {
	IDs        []string          `json:"ids"`
	ChainURL   string            `json:"chain_url"`
	ChainID    string            `json:"chain_id"`
	ChainType  string            `json:"chain_type"`
	StartBlock uint64            `json:"starting_block"`
	EndBlock   uint64            `json:"end_block"`
	Addresses  []string          `json:"addresses"`
	Standards  map[string]string `json:"standards"`
	EventTypes []string          `json:"event_types"`
}

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

func (h *Handler) Handle(ctx context.Context, input Input) (interface{}, error) {
	h.log.Debug().
		Uint64("start_block", input.StartBlock).
		Uint64("end_block", input.EndBlock).
		Strs("events", input.EventTypes).
		Strs("contracts", input.Addresses).
		Msg("processing job")

	network, err := web3.New(ctx, input.ChainURL)
	if err != nil {
		return nil, fmt.Errorf("could not create web3 client: %w", err)
	}
	defer network.Close()

	rawLogs, err := network.BlockEvents(ctx, input.StartBlock, input.EndBlock, input.EventTypes, input.Addresses)
	if err != nil {
		return nil, fmt.Errorf("could not get block events: %w", err)
	}

	logs := make([]log.Log, 0, len(rawLogs))
	for _, rawLog := range rawLogs {
		parser, err := h.getParser(network, rawLog.EventType)
		if err != nil {
			return nil, fmt.Errorf("could not get parser: %w", err)
		}

		log, err := parser.ParseRawLog(rawLog, input.Standards)
		if err != nil {
			return nil, fmt.Errorf("could not parse raw event: %w", err)
		}

		logs = append(logs, *log)
	}

	return logs, nil
}

func (h *Handler) getParser(network networks.Network, eventType string) (parsing.Parser, error) {
	parsers, err := h.initializer(network)
	if err != nil {
		return nil, fmt.Errorf("could not initialize parsers: %w", err)
	}

	for _, parser := range parsers {
		if strings.ToLower(parser.Type()) == strings.ToLower(eventType) {
			return parser, nil
		}
	}

	return nil, errParserNotFound
}
