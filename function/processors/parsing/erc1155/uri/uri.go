package uri

import (
	"bytes"
	"fmt"

	"github.com/NFT-com/indexer/function/processors/parsing/erc1155"
	"github.com/NFT-com/indexer/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
	return uriType
}

func (p *Parser) ParseRawLog(raw log.RawLog, standards map[string]string) ([]log.Log, error) {
	if len(raw.IndexData) < defaultIndexDataLen {
		return nil, fmt.Errorf("unexpected index data length (have: %d, want: %d)", len(raw.IndexData), defaultIndexDataLen)
	}

	var (
		id = common.HexToHash(raw.IndexData[0]).Big()
	)

	data := make(map[string]interface{})
	err := p.abi.UnpackIntoMap(data, eventName, raw.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack event: %w", err)
	}

	uri, ok := data[uriFieldName].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse uri: uri is empty or not a string value")
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
		Type:            log.URI,
		URI:             uri,
		EmittedAt:       raw.EmittedAt,
	}

	contractCollectionID, _ := erc1155.ExtractIDs(id.Bytes())

	// If it found the contract collection id, set it to the correct amount,
	// otherwise it was the contract collection 0.
	if contractCollectionID != nil {
		l.ContractCollectionID = contractCollectionID.String()
	} else {
		l.ContractCollectionID = "0"
	}

	return []log.Log{l}, nil
}
