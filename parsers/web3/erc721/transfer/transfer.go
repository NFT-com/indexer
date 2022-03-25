package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/log"
)

type Parser struct {
}

func NewParser() *Parser {
	p := Parser{}

	return &p
}

func (p *Parser) ParseRawEvent(rawLog log.RawLog) (*log.Log, error) {
	if len(rawLog.IndexData) < 3 {
		return nil, fmt.Errorf("could not parse raw log: index data lenght is less than 3")
	}

	var (
		fromAddress = common.HexToAddress(rawLog.IndexData[0]).String()
		toAddress   = common.HexToAddress(rawLog.IndexData[1]).String()
		nftID       = common.HexToHash(rawLog.IndexData[2]).Big().String()
	)

	l := log.Log{
		ID:              rawLog.ID,
		ChainID:         rawLog.ChainID,
		Contract:        rawLog.Address,
		Block:           rawLog.BlockNumber,
		TransactionHash: rawLog.TransactionHash,
		NftID:           nftID,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		EmittedAt:       rawLog.EmittedAt,
	}

	switch {
	case rawLog.IndexData[0] == zeroValueHash:
		l.Type = log.Mint
	case rawLog.IndexData[1] == zeroValueHash:
		l.Type = log.Burn
	default:
		l.Type = log.Transfer
	}

	return &l, nil
}
