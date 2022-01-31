package main

import (
	"context"
	"github.com/rs/zerolog"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/store/mock"
)

const (
	EnvVarLogLevel = "LOG_LEVEL"

	transferSingleEventName = "TransferSingle"
	transferBatchEventName  = "TransferBatch"
	uriEventName            = "URI"

	fromKeyword   = "from"
	toKeyword     = "to"
	idKeyword     = "id"
	idsKeyword    = "ids"
	valueKeyword  = "value"
	valuesKeyword = "values"
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

	data, err := abiEvent.Inputs.Unpack(e.Data)
	if err != nil {
		return err
	}

	switch abiEvent.Name {
	case transferSingleEventName:
		return h.handleSingleEvent(ctx, transferSingleEventName, e, data)
	case transferBatchEventName:
		return h.handleBatchEvent(ctx, transferBatchEventName, e, data)
	case uriEventName:
		return h.handleURIEvent(ctx, e, data)
	default:
		// We only care about the above events, for now other event is not worth unpack or saving
		return nil
	}
}

func (h *Handler) handleSingleEvent(ctx context.Context, name string, e *event.Event, data []interface{}) error {
	var (
		// We don't care about the operator for now, so just skipping that indexed field and skipping it.
		from  = e.IndexedData[2]
		to    = e.IndexedData[3]
		id    = *abi.ConvertType(data[0], new(*big.Int)).(**big.Int)
		value = *abi.ConvertType(data[1], new(*big.Int)).(**big.Int)
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
			fromKeyword:  from.Hex(),
			toKeyword:    to.Hex(),
			idKeyword:    id.String(),
			valueKeyword: value.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	if value.Cmp(big.NewInt(1)) != 0 {
		// We don't care about fungible tokens, so ignore it for now.
		return nil
	}

	switch {
	case from == common.HexToHash(""):
		return h.handleMintEvent(ctx, id, to, e)
	case to == common.HexToHash(""):
		return h.handleBurnEvent(ctx, id, e)
	default:
		return h.handleTransferEvent(ctx, id, to, e)
	}
}

func (h *Handler) handleBatchEvent(ctx context.Context, _ string, e *event.Event, data []interface{}) error {
	var (
		// We don't care about the operator for now, so just skipping that indexed field and skipping it.
		from   = e.IndexedData[2]
		to     = e.IndexedData[3]
		ids    = abi.ConvertType(data[0], make([]*big.Int, 0)).([]*big.Int) // FIXME: Test if this tests
		values = abi.ConvertType(data[1], make([]*big.Int, 0)).([]*big.Int)
	)

	// FIXME STORE EVENT

	for i, id := range ids {
		if values[i].Cmp(big.NewInt(1)) != 0 {
			// We don't care about fungible tokens, so ignore it for now.
			continue
		}

		switch {
		case from == common.HexToHash(""):
			return h.handleMintEvent(ctx, id, to, e)
		case to == common.HexToHash(""):
			return h.handleBurnEvent(ctx, id, e)
		default:
			return h.handleTransferEvent(ctx, id, to, e)
		}
	}

	return nil
}

func (h *Handler) handleURIEvent(ctx context.Context, e *event.Event, data []interface{}) error {
	log.Println("URI", data)
	// FIXME

	newURI := ""
	if err := h.store.UpdateContractURI(ctx, e.Network, e.Chain, e.Address.Hex(), newURI); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleMintEvent(ctx context.Context, id *big.Int, to common.Hash, e *event.Event) error {
	storeNFT := nft.NFT{
		ID:       id.String(),
		Network:  e.Network,
		Chain:    e.Chain,
		Contract: e.Address.Hex(),
		Owner:    to.Hex(),
	}

	if err := h.store.SaveNFT(ctx, &storeNFT); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleBurnEvent(ctx context.Context, id *big.Int, e *event.Event) error {
	if err := h.store.BurnNFT(ctx, e.Network, e.Chain, e.Address.Hex(), id.String()); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleTransferEvent(ctx context.Context, id *big.Int, to common.Hash, e *event.Event) error {
	if err := h.store.UpdateNFTOwner(ctx, e.Network, e.Chain, e.Address.Hex(), id.String(), to.Hex()); err != nil {
		return err
	}

	return nil
}
