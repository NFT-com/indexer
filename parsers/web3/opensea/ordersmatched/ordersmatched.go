package main

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/networks"
)

type Parser struct {
	client networks.Network
	abi    abi.ABI
}

func NewParser(client networks.Network) (*Parser, error) {
	if client == nil {
		return nil, fmt.Errorf("invalid argment: network client")
	}

	parsedAbi, err := abi.JSON(bytes.NewBufferString(ABI))
	if err != nil {
		return nil, fmt.Errorf("could not parse abi: %w", err)
	}

	p := Parser{
		client: client,
		abi:    parsedAbi,
	}

	return &p, nil
}

func (p *Parser) ParseRawEvent(rawEvent events.RawEvent) (events.Event, error) {
	seller := rawEvent.IndexData[1]
	buyer := rawEvent.IndexData[2]

	order := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(order, EventName, rawEvent.Data)
	if err != nil {
		return events.Event{}, fmt.Errorf("could not unpack events: %w", err)
	}

	price, ok := order[PriceFieldName].(*big.Int)
	if !ok {
		return events.Event{}, fmt.Errorf("could not parse price: price is not a big.Int pointer")
	}

	m := events.Event{
		ID:          rawEvent.ID,
		Type:        events.TypeSell,
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
