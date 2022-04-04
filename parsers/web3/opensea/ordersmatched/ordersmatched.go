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
		return nil, fmt.Errorf("invalid argument: nil network client")
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

func (p *Parser) ParseRawLog(raw log.RawLog) (*log.Log, error) {
	if len(raw.IndexData) != defaultIndexDataLen {
		return nil, fmt.Errorf("unexpected index data length (have: %v, want: %v)", len(raw.IndexData), defaultIndexDataLen)
	}

	var (
		seller = raw.IndexData[0]
		buyer  = raw.IndexData[1]
	)

	data := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(data, eventName, raw.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log data: %w", err)
	}

	price, ok := data[priceFieldName].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("could not parse price: price is empty or not a big.Int pointer")
	}

	l := log.Log{
		ID:              raw.ID,
		Type:            log.Sale,
		ChainID:         raw.ChainID,
		Block:           raw.BlockNumber,
		Index:           raw.Index,
		TransactionHash: raw.TransactionHash,
		Contract:        raw.Address,
		FromAddress:     common.HexToAddress(seller).String(),
		ToAddress:       common.HexToAddress(buyer).String(),
		Price:           price.String(),
		EmittedAt:       raw.EmittedAt,
	}

	return &l, nil
}
