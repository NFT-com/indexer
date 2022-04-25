package action

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function/processors/action"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/networks/web3"
)

var (
	errParserNotFound = errors.New("parser not found")
)

// Initializer initializes the processors to use with the network client.
type Initializer func(client networks.Network) ([]action.Processor, error)

type Handler struct {
	log         zerolog.Logger
	initializer Initializer
}

func NewHandler(log zerolog.Logger, initializer Initializer) *Handler {
	h := Handler{
		log:         log,
		initializer: initializer,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, job jobs.Action) (*chain.NFT, error) {
	h.log.Debug().
		Str("block", job.BlockNumber).
		Str("event", job.Event).
		Str("contract", job.Address).
		Msg("processing job")

	network, err := web3.New(ctx, job.ChainURL)
	if err != nil {
		return nil, fmt.Errorf("could not create web3 client: %w", err)
	}
	defer network.Close()

	processor, err := h.getProcessor(network, job.Type, job.Standard)
	if err != nil {
		return nil, fmt.Errorf("could not get processor: %w", err)
	}

	nft, err := processor.Process(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("could not process action: %w", err)
	}

	return nft, nil
}

func (h *Handler) getProcessor(network networks.Network, actionType, standard string) (action.Processor, error) {
	processors, err := h.initializer(network)
	if err != nil {
		return nil, fmt.Errorf("could not initialize parsers: %w", err)
	}

	for _, processor := range processors {
		if strings.ToLower(processor.Type()) == strings.ToLower(actionType) &&
			strings.ToLower(processor.Standard()) == strings.ToLower(standard) {
			return processor, nil
		}
	}

	return nil, errParserNotFound
}
