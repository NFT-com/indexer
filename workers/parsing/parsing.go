package parsing

import (
	"context"
	"github.com/NFT-com/indexer/events/mint"
	"github.com/NFT-com/indexer/events/update"
	"github.com/rs/zerolog"
	"log"

	"github.com/NFT-com/indexer/job"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/parsers/erc721/transfer"
)

type Handler struct {
	log   zerolog.Logger
	store Store
}

func NewHandler(log zerolog.Logger, store Store) *Handler {
	h := Handler{
		log:   log.With().Str("component", "parsing_handler").Logger(),
		store: store,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, parsingJob job.Parsing) error {
	network, err := web3.NewWeb3(ctx, parsingJob.ChainURL)
	if err != nil {
		return err
	}
	defer network.Close()

	parser := transfer.NewParser()

	rawEvents, err := network.BlockEvents(ctx, parsingJob.BlockNumber, parsingJob.EventType, parsingJob.Address)
	if err != nil {
		return err
	}

	for _, rawEvent := range rawEvents {
		err = h.store.InsertRawEvent(rawEvent)
		if err != nil {
			h.log.Error().Err(err).Msg("failed to store raw events")
			continue
		}

		parsedEvent, err := parser.ParseRawEvent(rawEvent)
		if err != nil {
			h.log.Error().Err(err).Msg("failed to parse raw events")
			continue
		}

		switch event := parsedEvent.(type) {
		case mint.Mint:
			log.Println(h.store.InsertNewNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, event.ToAddress))
		case update.Update:
			log.Println(h.store.InsertNewNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, "0x0"))
			log.Println(h.store.UpdateNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, event.ToAddress))
		}
	}

	return nil
}
