package parsing

import (
	"context"
	"errors"
	"github.com/NFT-com/indexer/events/mint"
	"github.com/NFT-com/indexer/events/update"
	"github.com/NFT-com/indexer/job"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/parsers/erc721/transfer"
	"github.com/NFT-com/indexer/service/postgres"
	"github.com/rs/zerolog"
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
	log := h.log.With().
		Str("block", parsingJob.BlockNumber).
		Str("event", parsingJob.EventType).
		Str("contract", parsingJob.Address).
		Logger()

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
		if err != nil && !errors.Is(err, postgres.ErrAlreadyExists) {
			log.Error().Err(err).Msg("failed to store raw events")
			return err
		}

		parsedEvent, err := parser.ParseRawEvent(rawEvent)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse raw events")
			return err
		}

		switch event := parsedEvent.(type) {
		case mint.Mint:
			err = h.store.InsertNewNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, event.ToAddress)
			if err != nil {
				log.Error().Err(err).Msg("failed to insert new nft in mint event")
				return err
			}
		case update.Update:
			err = h.store.UpdateNFT(event.NetworkID, event.ChainID, event.Contract, event.NftID, event.ToAddress)
			if err != nil {
				log.Error().Err(err).Msg("failed to update nft on update event")
				return err
			}
		}
	}

	return nil
}
