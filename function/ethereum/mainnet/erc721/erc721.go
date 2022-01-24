package erc721

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
	"github.com/NFT-com/indexer/store"
)

const (
	Name              = "Ethereum-mainnet-erc721"
	transferEventName = "Transfer"
)

type Handler struct {
	store store.Storer
}

func New(store store.Storer) *Handler {
	h := Handler{
		store: store,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, e *event.Event) error {
	jsonABI, err := h.store.GetContractABI(ctx, e.Network, e.Chain, e.Address.Hex())
	if err != nil {
		return err
	}

	parsedABI, err := abi.JSON(strings.NewReader(jsonABI))
	if err != nil {
		return err
	}

	abiEvent, err := parsedABI.EventByID(e.Topic)
	if err != nil {
		return err
	}

	if abiEvent.Name != transferEventName {
		// We only care about transfer events, for now approve event is not worth unpack or saving
		return nil
	}

	if len(e.IndexedData) != 3 {
		// This needs to have from, to and id fields
		return nil
	}

	var (
		from = e.IndexedData[0]
		to   = e.IndexedData[1]
		id   = e.IndexedData[2].Big()
	)

	parsedEvent := event.ParsedEVent{
		ID:              e.ID,
		Network:         e.Network,
		Chain:           e.Chain,
		Block:           e.Block,
		TransactionHash: e.TransactionHash.Hex(),
		Address:         e.Address.Hex(),
		Type:            abiEvent.Name,
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	if from == common.HexToHash("") {
		// FIXME: GET METADATA

		storeNFT := nft.NFT{
			ID:       id.String(),
			Network:  e.Network,
			Chain:    e.Chain,
			Contract: e.Address.Hex(),
			Owner:    to.Hex(),
			// FIXME: DATA
		}

		if err := h.store.SaveNFT(ctx, &storeNFT); err != nil {
			return err
		}
	}

	return nil
}
