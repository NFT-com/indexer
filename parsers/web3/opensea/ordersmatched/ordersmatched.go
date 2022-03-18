package main

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/networks"
)

type Parser struct {
	client networks.Network
	abi    abi.ABI
}

func NewParser(client networks.Network) (*Parser, error) {
	parsedAbi, err := abi.JSON(bytes.NewBufferString(ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %w", err)
	}

	p := Parser{
		client: client,
		abi:    parsedAbi,
	}

	return &p, nil
}

func (p *Parser) ParseRawEvent(rawEvent event.RawEvent) (event.Event, error) {
	seller := rawEvent.IndexData[1]
	buyer := rawEvent.IndexData[2]

	order := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(order, EventName, rawEvent.Data)
	if err != nil {
		return event.Event{}, fmt.Errorf("failed to unpack event: %w", err)
	}

	price, ok := order[PriceFieldName].(*big.Int)
	if !ok {
		return event.Event{}, fmt.Errorf("failed to parse price: price is not a big.Int pointer")
	}

	m := event.Event{
		ID:          rawEvent.ID,
		Type:        event.TypeSell,
		ChainID:     rawEvent.ChainID,
		NetworkID:   rawEvent.NetworkID,
		Contract:    rawEvent.Address,
		FromAddress: common.HexToAddress(seller).String(),
		ToAddress:   common.HexToAddress(buyer).String(),
		Price:       price.String(),
		EmittedAt:   rawEvent.EmittedAt,
	}

	return m, nil
}
