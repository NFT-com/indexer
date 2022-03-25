package main

import (
	"fmt"

	"github.com/NFT-com/indexer/event"
	"github.com/ethereum/go-ethereum/common"
)

const (
	zeroValueHash = "0x0000000000000000000000000000000000000000"
)

type Parser struct {
}

func NewParser() *Parser {
	p := Parser{}

	return &p
}

func (p *Parser) ParseRawEvent(rawEvent event.RawEvent) (*event.Event, error) {
	if len(rawEvent.IndexData) < 3 {
		return nil, fmt.Errorf("could not parse raw event: index data lenght is less than 3")
	}

	var (
		fromAddress = common.HexToAddress(rawEvent.IndexData[0]).String()
		toAddress   = common.HexToAddress(rawEvent.IndexData[1]).String()
		nftID       = common.HexToHash(rawEvent.IndexData[2]).Big().String()
	)

	m := event.Event{
		ID:          rawEvent.ID,
		ChainID:     rawEvent.ChainID,
		NetworkID:   rawEvent.NetworkID,
		NftID:       nftID,
		Contract:    rawEvent.Address,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		EmittedAt:   rawEvent.EmittedAt,
	}

	switch {
	case rawEvent.IndexData[0] == zeroValueHash:
		m.Type = event.Mint
	case rawEvent.IndexData[1] == zeroValueHash:
		m.Type = event.Burn
	default:
		m.Type = event.Transfer
	}

	return &m, nil
}
