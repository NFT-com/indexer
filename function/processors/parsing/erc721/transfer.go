package erc721

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

func (p *Parser) Type() string {
	return transferType
}

func (p *Parser) ParseRawLog(raw log.RawLog, standards map[string]string) ([]log.Log, error) {
	if len(raw.IndexData) != defaultIndexDataLen {
		return nil, fmt.Errorf("unexpected index data length (have: %d, want: %d)", len(raw.IndexData), defaultIndexDataLen)
	}

	var (
		fromAddress = common.HexToAddress(raw.IndexData[0]).String()
		toAddress   = common.HexToAddress(raw.IndexData[1]).String()
		nftID       = common.HexToHash(raw.IndexData[2]).Big().String()
	)

	l := log.Log{
		ID:              raw.ID,
		ChainID:         raw.ChainID,
		Contract:        raw.Address,
		Block:           raw.BlockNumber,
		Standard:        standards[raw.EventType],
		Event:           raw.EventType,
		Index:           raw.Index,
		TransactionHash: raw.TransactionHash,
		NeedsActionJob:  true,
		NftID:           nftID,
		Amount:          1,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		EmittedAt:       raw.EmittedAt,
	}

	switch zeroValueAddress {
	case fromAddress:
		l.Type = log.Mint
		l.ActionType = log.Addition
	case toAddress:
		l.Type = log.Burn
		l.ActionType = log.OwnerChange
	default:
		l.Type = log.Transfer
		l.ActionType = log.OwnerChange
	}

	return []log.Log{l}, nil
}
