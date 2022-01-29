package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/nft"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/store/mock"
)

const (
	EnvVarNodeURL = "NODE_URL"

	birthEventName    = "Birth"
	transferEventName = "Transfer"

	fromKeyword     = "from"
	toKeyword       = "to"
	idKeyword       = "id"
	ownerKeyword    = "owner"
	kittyIDKeyword  = "kitty_id"
	matronIDKeyword = "matron_id"
	sireIDKeyword   = "sire_id"
	genesKeyword    = "genes"
)

func main() {
	val, ok := os.LookupEnv(EnvVarNodeURL)
	if !ok {
		log.Fatalln("missing environment variable")
	}

	store := mock.New()

	client, err := ethclient.Dial(val)
	if err != nil {
		log.Fatalln(err)
	}

	handler := New(store, client)

	lambda.Start(handler.Handle)
}

type Handler struct {
	store  store.Storer
	client *ethclient.Client
}

func New(store store.Storer, client *ethclient.Client) *Handler {
	h := Handler{
		store:  store,
		client: client,
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
	case birthEventName:
		return h.handleBirthEvent(ctx, birthEventName, e, data)
	case transferEventName:
		return h.handleTransferEvent(ctx, transferEventName, e, data)
	default:
		return nil
	}
}

func (h *Handler) handleBirthEvent(ctx context.Context, name string, e *event.Event, data []interface{}) error {
	var (
		owner    = abi.ConvertType(data[0], new(common.Address)).(*common.Address)
		kittyID  = abi.ConvertType(data[1], new(big.Int)).(*big.Int)
		matronID = abi.ConvertType(data[2], new(big.Int)).(*big.Int)
		sireID   = abi.ConvertType(data[3], new(big.Int)).(*big.Int)
		genes    = abi.ConvertType(data[4], new(big.Int)).(*big.Int)
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
			ownerKeyword:    owner.Hex(),
			kittyIDKeyword:  kittyID.String(),
			matronIDKeyword: matronID.String(),
			sireIDKeyword:   sireID.String(),
			genesKeyword:    genes.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	storeNFT := nft.NFT{
		ID:       kittyID.String(),
		Network:  e.Network,
		Chain:    e.Chain,
		Contract: e.Address.Hex(),
		Owner:    owner.Hex(),
		Data: map[string]interface{}{
			matronIDKeyword: matronID.String(),
			sireIDKeyword:   sireID.String(),
			genesKeyword:    genes.String(),
		},
	}

	if err := h.store.SaveNFT(ctx, &storeNFT); err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleTransferEvent(ctx context.Context, name string, e *event.Event, data []interface{}) error {
	var (
		from = *abi.ConvertType(data[0], new(common.Hash)).(*common.Hash)
		to   = *abi.ConvertType(data[1], new(common.Hash)).(*common.Hash)
		id   = abi.ConvertType(data[2], new(big.Int)).(*big.Int)
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
			fromKeyword: from.Hex(),
			toKeyword:   to.Hex(),
			idKeyword:   id.String(),
		},
	}

	if err := h.store.SaveEvent(ctx, &parsedEvent); err != nil {
		return err
	}

	switch {
	case from == common.HexToHash(""):
		// Already handled by the birth event
		return nil
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
