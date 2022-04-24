package addition

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function/processors/addition/erc721metadata"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/networks/web3"
)

type Handler struct {
	log zerolog.Logger
}

func NewHandler(log zerolog.Logger) *Handler {
	h := Handler{
		log: log,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, job jobs.Addition) (*chain.NFT, error) {
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

	processor, err := erc721metadata.NewProcessor(h.log, network)
	if err != nil {
		return nil, fmt.Errorf("could not create processor: %w", err)
	}

	nft, err := processor.Process(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("could not process additon: %w", err)
	}

	return nft, nil
}
