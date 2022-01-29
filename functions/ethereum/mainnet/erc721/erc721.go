package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum"
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

	uriMethod         = "tokenURI"
	transferEventName = "Transfer"

	fromKeyword = "from"
	toKeyword   = "to"
	idKeyword   = "id"
	uriKeyword  = "uri"
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

	parsedEvent := event.ParsedEvent{
		ID:              e.ID,
		Network:         e.Network,
		Chain:           e.Chain,
		Block:           e.Block,
		TransactionHash: e.TransactionHash.Hex(),
		Address:         e.Address.Hex(),
		Type:            abiEvent.Name,
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
		return h.handleMintEvent(ctx, parsedABI, id, to, e)
	case to == common.HexToHash(""):
		return h.handleBurnEvent(ctx, id, e)
	default:
		return h.handleTransferEvent(ctx, id, to, e)
	}
}

func (h *Handler) handleMintEvent(ctx context.Context, parsedABI abi.ABI, id *big.Int, to common.Hash, e *event.Event) error {
	input, err := parsedABI.Pack(uriMethod, id)
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{To: &e.Address, Data: input}
	data, err := h.client.CallContract(ctx, msg, nil)
	if err != nil {
		return err
	}

	unpackedData, err := parsedABI.Unpack(uriMethod, data)
	if err != nil {
		return err
	}

	if len(unpackedData) == 0 {
		return nil // FIXME
	}

	uri, ok := unpackedData[0].(string)
	if !ok {
		return nil // FIXME
	}

	storeNFT := nft.NFT{
		ID:       id.String(),
		Network:  e.Network,
		Chain:    e.Chain,
		Contract: e.Address.Hex(),
		Owner:    to.Hex(),
		Data: map[string]interface{}{
			uriKeyword: uri,
		},
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
