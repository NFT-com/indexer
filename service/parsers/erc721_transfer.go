package parsers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

func ERC721Transfer(log types.Log) (*events.Transfer, error) {

	if len(log.Topics) < 3 {
		return nil, fmt.Errorf("invalid topic length have (%d) want >= (%d)", len(log.Topics), 3)
	}

	transfer := events.Transfer{
		ID: logID(log),
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
