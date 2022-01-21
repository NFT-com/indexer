package erc1155

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/store"
)

const (
	Name                    = "Ethereum-mainnet-erc1155"
	transferSingleEventName = "TransferSingle"
	transferBatchEventName  = "TransferBatch"
	uriEventName            = "URI"
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

	switch abiEvent.Name {
	case transferSingleEventName:
		// FIXME: Do the transfer single
	case transferBatchEventName:
		// FIXME: Do the transfer batch
	case uriEventName:
		// FIXME: Do the uri
	default:
		// We only care about the above events, for now other event is not worth unpack or saving
		return nil
	}

	data, err := abiEvent.Inputs.Unpack(e.Data)
	if err != nil {
		return err
	}

	fmt.Println(data)

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

	return nil
}