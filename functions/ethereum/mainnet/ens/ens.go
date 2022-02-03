package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/store/mock"
)

const (
	EnvVarLogLevel = "LOG_LEVEL"

	nameMigratedEventName = "NameMigrated"
	nameRegisterEventName = "NameRegistered"
	nameRenewedEventName  = "NameRenewed"
	transferEventName     = "Transfer"

	fromKeyword    = "from"
	toKeyword      = "to"
	idKeyword      = "id"
	ownerKeyword   = "owner"
	expiresKeyword = "expires"
)

func main() {
	logLevel, ok := os.LookupEnv(EnvVarLogLevel)
	if !ok {
		logLevel = "info"
	}

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Fatalln("failed to parse log level", err)
	}
	logger = logger.Level(level)

	store := mock.New(logger)
	handler := New(store)

	lambda.Start(handler.Handle)
}

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
	case nameMigratedEventName:
		return h.handleNameMigratedEvent(ctx, nameMigratedEventName, e, abiEvent)
	case nameRegisterEventName:
		return h.handleNameRegisterEvent(ctx, nameRegisterEventName, e, abiEvent)
	case nameRenewedEventName:
		return h.handleNameRenewedEvent(ctx, nameRenewedEventName, e, abiEvent)
	case transferEventName:
		return h.handleTransferEvent(ctx, transferEventName, e)
	default:
		return nil
	}
}

func (h *Handler) handleNameMigratedEvent(ctx context.Context, name string, e *event.Event, abiEvent *abi.Event) error {
	data, err := abiEvent.Inputs.Unpack(e.Data)
	if err != nil {
		return err
	}

	var (
		id      = e.IndexedData[0].Big()
		owner   = e.IndexedData[1]
		expires = abi.ConvertType(data[0], new(big.Int)).(*big.Int)
	)

	parsedEvent := event.ParsedEvent{
		ID:              e.ID,
		Network:         e.Network,
		Chain:           e.Chain,
		Block:           e.Block,
		TransactionHash: e.TransactionHash.Hex(),
		Address:         e.Address.Hex(),
		Type:            name,
		Data: map[string]interface{}{
			idKeyword:      id.String(),
			ownerKeyword:   owner.Hex(),
			expiresKeyword: expires.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	if err := h.store.UpdateNFTOwner(ctx, e.Network, e.Chain, e.Address.Hex(), id.String(), owner.Hex()); err != nil {
		return err
	}

	newData := map[string]interface{}{
		expiresKeyword: expires.String(),
	}
	if err := h.store.UpdateNFTMetadata(ctx, e.Network, e.Chain, e.Address.Hex(), id.String(), newData); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleNameRegisterEvent(ctx context.Context, name string, e *event.Event, abiEvent *abi.Event) error {
	data, err := abiEvent.Inputs.Unpack(e.Data)
	if err != nil {
		return err
	}

	var (
		id      = e.IndexedData[0].Big()
		owner   = e.IndexedData[1]
		expires = abi.ConvertType(data[0], new(big.Int)).(*big.Int)
	)

	parsedEvent := event.ParsedEvent{
		ID:              e.ID,
		Network:         e.Network,
		Chain:           e.Chain,
		Block:           e.Block,
		TransactionHash: e.TransactionHash.Hex(),
		Address:         e.Address.Hex(),
		Type:            name,
		Data: map[string]interface{}{
			idKeyword:      id.String(),
			ownerKeyword:   owner.Hex(),
			expiresKeyword: expires.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	storeNFT := nft.NFT{
		ID:       id.String(),
		Network:  e.Network,
		Chain:    e.Chain,
		Contract: e.Address.Hex(),
		Owner:    owner.Hex(),
		Data: map[string]interface{}{
			expiresKeyword: expires.String(),
		},
	}

	if err := h.store.SaveNFT(ctx, &storeNFT); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleNameRenewedEvent(ctx context.Context, name string, e *event.Event, abiEvent *abi.Event) error {
	data, err := abiEvent.Inputs.Unpack(e.Data)
	if err != nil {
		return err
	}

	var (
		id      = e.IndexedData[0].Big()
		expires = abi.ConvertType(data[0], new(big.Int)).(*big.Int)
	)

	parsedEvent := event.ParsedEvent{
		ID:              e.ID,
		Network:         e.Network,
		Chain:           e.Chain,
		Block:           e.Block,
		TransactionHash: e.TransactionHash.Hex(),
		Address:         e.Address.Hex(),
		Type:            name,
		Data: map[string]interface{}{
			idKeyword:      id.String(),
			expiresKeyword: expires.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	newData := map[string]interface{}{
		expiresKeyword: expires.String(),
	}
	if err := h.store.UpdateNFTMetadata(ctx, e.Network, e.Chain, e.Address.Hex(), id.String(), newData); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleTransferEvent(ctx context.Context, name string, e *event.Event) error {
	var (
		from = e.IndexedData[0]
		to   = e.IndexedData[1]
		id   = e.IndexedData[2].Big()
	)

	if from == common.HexToHash("") {
		// Already handled by the name registered event
		return nil
	}

	parsedEvent := event.ParsedEvent{
		ID:              e.ID,
		Network:         e.Network,
		Chain:           e.Chain,
		Block:           e.Block,
		TransactionHash: e.TransactionHash.Hex(),
		Address:         e.Address.Hex(),
		Type:            name,
		Data: map[string]interface{}{
			fromKeyword: from.Hex(),
			toKeyword:   to.Hex(),
			idKeyword:   id.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	switch {
	case to == common.HexToHash(""):
		return h.handleBurnEvent(ctx, id, e)
	default:
		return h.handleSwitchEvent(ctx, id, to, e)
	}
}

func (h *Handler) handleBurnEvent(ctx context.Context, id *big.Int, e *event.Event) error {
	if err := h.store.BurnNFT(ctx, e.Network, e.Chain, e.Address.Hex(), id.String()); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleSwitchEvent(ctx context.Context, id *big.Int, to common.Hash, e *event.Event) error {
	if err := h.store.UpdateNFTOwner(ctx, e.Network, e.Chain, e.Address.Hex(), id.String(), to.Hex()); err != nil {
		return err
	}

	return nil
}
