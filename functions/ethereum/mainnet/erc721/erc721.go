package main

import (
	"context"
	"log"
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

	transferEventName = "Transfer"

	fromKeyword = "from"
	toKeyword   = "to"
	idKeyword   = "id"
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

	parsedEvent := event.ParsedEVent{
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
