package parsers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/id"
	"github.com/NFT-com/indexer/models/jobs"
)

const (
	erc20Transfer = "Transfer"
	erc20Value    = "value"
)

// ERC20Transfer takes a log entry and parses it as an ERC20 transfer.
func ERC20Transfer(log types.Log) (*events.Transfer, error) {

	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("invalid number of topics (want: %d, have: %d)", 3, len(log.Topics))
	}

	fields := make(map[string]interface{})
	err := abis.ERC20.UnpackIntoMap(fields, erc20Transfer, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack fields: %w", err)
	}

	if len(fields) != 1 {
		return nil, fmt.Errorf("invalid number of fields (want: %d, have: %d)", 1, len(fields))
	}

	valueField, ok := fields[erc20Value]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", erc20Value)
	}

	value, ok := valueField.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type (field: %s,  have: %T)", erc20Value, valueField)
	}

	transfer := events.Transfer{
		ID: id.Log(log),
		// ChainID set after parsing
		TokenStandard:     jobs.StandardERC20,
		CollectionAddress: log.Address.Hex(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		SenderAddress:     common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
		ReceiverAddress:   common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		TokenCount:        value.String(),
		// EmittedAt set after parsing
	}

	return &transfer, nil
}
