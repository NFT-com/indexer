package parsers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/id"
	"github.com/NFT-com/indexer/models/jobs"
)

const (
	erc721Transfer = "Transfer"
)

func ERC721Transfer(log types.Log) (*events.Transfer, error) {

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid number of topics (want: %d, have: %d)", 4, len(log.Topics))
	}

	fields := make(map[string]interface{})
	err := abis.ERC721.UnpackIntoMap(fields, erc721Transfer, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack fields: %w", err)
	}

	if len(fields) != 0 {
		return nil, fmt.Errorf("invalid number of fields (want: %d, have: %d)", 0, len(fields))
	}

	transfer := events.Transfer{
		ID: id.Log(log),
		// ChainID set after parsing
		TokenStandard:     jobs.StandardERC721,
		CollectionAddress: log.Address.Hex(),
		TokenID:           log.Topics[3].Big().String(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		SenderAddress:     common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
		ReceiverAddress:   common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		TokenCount:        "1",
		// EmittedAt set after parsing
	}

	return &transfer, nil
}
