package main

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/log"
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

func (p *Parser) ParseRawLog(rawLog log.RawLog) (*log.Log, error) {
	if len(rawLog.IndexData) < 2 {
		return nil, fmt.Errorf("could not parse raw log: index data lenght is less than 2")
	}

	var (
		seller = rawLog.IndexData[0]
		buyer  = rawLog.IndexData[1]
	)

	data := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(data, eventName, rawLog.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log data: %w", err)
	}

	price, ok := data[priceFieldName].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("could not parse price: price is empty or not a big.Int pointer")
	}

	l := log.Log{
		ID:              rawLog.ID,
		Type:            log.Sale,
		ChainID:         rawLog.ChainID,
		Block:           rawLog.BlockNumber,
		TransactionHash: rawLog.TransactionHash,
		Contract:        rawLog.Address,
		FromAddress:     common.HexToAddress(seller).String(),
		ToAddress:       common.HexToAddress(buyer).String(),
		Price:           price.String(),
		EmittedAt:       rawLog.EmittedAt,
	}

	return &l, nil
}
