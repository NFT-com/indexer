package parsers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

const (
	eventTransferSingle = "TransferSingle"
	fieldID             = "id"
)

func ERC1155Transfer(log types.Log) (*events.Transfer, error) {

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid topic lenght have (%d) want (%d)", len(log.Topics), 4)
	}

	fields := make(map[string]interface{})
	err := abis.ERC1155.UnpackIntoMap(fields, eventTransferSingle, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	tokenID, ok := fields[fieldID].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for %q field (%T)", fieldID, fields[fieldID])
	}
	count, ok := fields[fieldValue].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for %q field (%T)", fieldValue, fields[fieldValue])
	}

	transfer := events.Transfer{
		ID: logID(log),
		// ChainID set after parsing
		TokenStandard:     jobs.StandardERC1155,
		CollectionAddress: log.Address.Hex(),
		TokenID:           tokenID.String(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		SenderAddress:     common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		ReceiverAddress:   common.BytesToAddress(log.Topics[3].Bytes()).Hex(),
		TokenCount:        count.String(),
		// EmittedAt set after parsing
	}

	return &transfer, nil
}
