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
	eventTransferBatch = "TransferBatch"
	fieldIDs           = "ids"
	fieldValues        = "values"
)

func ERC1155Batch(log types.Log) ([]*events.Transfer, error) {

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid topic lenght have (%d) want (%d)", len(log.Topics), 4)
	}

	fields := make(map[string]interface{})
	err := abis.ERC1155.UnpackIntoMap(fields, eventTransferBatch, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	tokenIDs, ok := fields[fieldIDs].([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for %q field (%T)", fieldIDs, fields[fieldIDs])
	}
	counts, ok := fields[fieldValues].([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for %q field (%T)", fieldValues, fields[fieldValues])
	}

	var transfers []*events.Transfer
	for i, tokenID := range tokenIDs {
		count := counts[i]

		data := make([]byte, 8+32+8+32)
		binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
		copy(data[8:40], log.TxHash[:])
		binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
		copy(data[48:80], tokenID.Bytes())
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
		transfers = append(transfers, &transfer)
	}

	return transfers, nil
}
