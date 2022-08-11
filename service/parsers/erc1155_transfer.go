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
	erc1155Transfer = "TransferSingle"
	erc1155ID       = "id"
	erc1155Value    = "value"
)

func ERC1155Transfer(log types.Log) (*events.Transfer, error) {

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid number of topics (want: %d, have: %d)", 4, len(log.Topics))
	}

	fields := make(map[string]interface{})
	err := abis.ERC1155.UnpackIntoMap(fields, erc1155Transfer, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid number of fields (want: %d, have: %d)", 2, len(fields))
	}

	fieldID, ok := fields[erc1155ID]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", erc1155ID)
	}

	fieldValue, ok := fields[erc1155Value]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", erc1155Value)
	}

	tokenID, ok := fieldID.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type (field: %s, want: %T, have: %T)", erc1155ID, &big.Int{}, fieldID)
	}

	value, ok := fieldValue.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type (field: %s, want: %T, have: %T)", erc1155Value, &big.Int{}, fieldValue)
	}

	transfer := events.Transfer{
		ID: id.Log(log),
		// ChainID set after parsing
		TokenStandard:     jobs.StandardERC1155,
		CollectionAddress: log.Address.Hex(),
		TokenID:           tokenID.String(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		SenderAddress:     common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		ReceiverAddress:   common.BytesToAddress(log.Topics[3].Bytes()).Hex(),
		TokenCount:        value.String(),
		// EmittedAt set after parsing
	}

	return &transfer, nil
}
