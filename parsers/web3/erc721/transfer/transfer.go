package main

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	ZeroValueHash = "0x0000000000000000000000000000000000000000"
)

type Parser struct {
}

func NewParser() *Parser {
	p := Parser{}

	return &p
}

func (p *Parser) ParseRawEvent(rawEvent events.RawEvent) (events.Event, error) {
	var (
		fromAddress = common.HexToAddress(rawEvent.IndexData[0]).String()
		toAddress   = common.HexToAddress(rawEvent.IndexData[1]).String()
		nftID       = common.HexToHash(rawEvent.IndexData[2]).Big().String()
	)

	m := events.Event{
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
	case rawEvent.IndexData[0] == ZeroValueHash:
		m.Type = events.TypeMint
	case rawEvent.IndexData[1] == ZeroValueHash:
		m.Type = events.TypeBurn
	default:
		m.Type = events.TypeTransfer
	}

	return m, nil
}
