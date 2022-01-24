package mainnet

import (
	"context"
	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/store"
)

const (
	Name = "Ethereum-mainnet"
)

type Handler struct {
	store      store.Storer
	dispatcher dispatch.Dispatcher
}

func New(store store.Storer, dispatcher dispatch.Dispatcher) *Handler {
	h := Handler{
		store:      store,
		dispatcher: dispatcher,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, e *event.Event) error {
	contractType, err := h.store.GetContractType(ctx, e.Network, e.Chain, e.Address.Hex())
	if err != nil {
		return err
	}

	newFunctionName := ""
	switch contractType {
	case "erc721":
		newFunctionName = function.Name(e.Network, e.Chain, "erc721")
	case "erc1155":
		newFunctionName = function.Name(e.Network, e.Chain, "erc1155")
	case "custom":
		newFunctionName = function.Name(e.Network, e.Chain, e.Address.Hex())
	}

	if err := h.dispatcher.Dispatch(newFunctionName, e); err != nil {
		return err
	}

	return nil
}
