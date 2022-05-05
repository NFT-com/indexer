package singletransfer

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/log"
)

type Parser struct {
	abi abi.ABI
}

func NewParser() (*Parser, error) {
	parsedABI, err := abi.JSON(bytes.NewBufferString(eventABI))
	if err != nil {
		return nil, fmt.Errorf("could not parse abi: %w", err)
	}

	p := Parser{
		abi: parsedABI,
	}

	return &p, nil
}

func (p *Parser) Type() string {
	return transferType
}

func (p *Parser) ParseRawLog(raw log.RawLog, standards map[string]string) ([]log.Log, error) {
	if len(raw.IndexData) != defaultIndexDataLen {
		return nil, fmt.Errorf("unexpected index data length (have: %d, want: %d)", len(raw.IndexData), defaultIndexDataLen)
	}

	var (
		// we don't care about the operator so just start on the index 1
		fromAddress = common.HexToAddress(raw.IndexData[1]).String()
		toAddress   = common.HexToAddress(raw.IndexData[2]).String()
	)

	data := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(data, eventName, raw.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack data into map: %w", err)
	}

	id, ok := data[idFieldName].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("could not parse id: id is empty or not a big.Int pointer")
	}

	value, ok := data[valueFieldName].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("could not parse value: value is empty or not a big.Int pointer")
	}

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
		NftID:           id.String(),
		Amount:          value.Uint64(),
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
