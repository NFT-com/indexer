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
	if client == nil {
		return nil, fmt.Errorf("invalid argment: network client")
	}

	parsedABI, err := abi.JSON(bytes.NewBufferString(eventABI))
	if err != nil {
		return nil, fmt.Errorf("could not parse abi: %w", err)
	}

	p := Parser{
		client: client,
		abi:    parsedABI,
	}

	return &p, nil
}

func (p *Parser) ParseRawEvent(rawEvent event.RawEvent) (*event.Event, error) {
	if len(rawEvent.IndexData) < 2 {
		return nil, fmt.Errorf("could not parse raw event: index data lenght is less than 2")
	}

	var (
		seller = rawEvent.IndexData[0]
		buyer  = rawEvent.IndexData[1]
	)

	order := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(order, eventName, rawEvent.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack event: %w", err)
	}

	price, ok := order[priceFieldName].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("could not parse price: price is empty or not a big.Int pointer")
	}

	m := event.Event{
		ID:          rawEvent.ID,
		Type:        event.Sale,
		ChainID:     rawEvent.ChainID,
		Contract:    rawEvent.Address,
		FromAddress: common.HexToAddress(seller).String(),
		ToAddress:   common.HexToAddress(buyer).String(),
		Price:       price.String(),
		EmittedAt:   rawEvent.EmittedAt,
	}

	return &m, nil
}
