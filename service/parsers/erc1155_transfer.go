package parsers

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

const (
	eventTransferSingle = "TransferSingle"
	fieldID             = "ids"
)

func ERC1155Transfer(log types.Log) (*events.Transfer, error) {

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

	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
	hash := sha3.Sum256(data)
	transferID := uuid.Must(uuid.FromBytes(hash[:16]))

	transfer := events.Transfer{
		ID: transferID.String(),
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
