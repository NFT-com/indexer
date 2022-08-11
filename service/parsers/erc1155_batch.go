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
	erc1155Batch  = "TransferBatch"
	erc1155IDs    = "ids"
	erc1155Values = "values"
)

func ERC1155Batch(log types.Log) ([]*events.Transfer, error) {

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid number of topics (want: %d, have: %d)", 4, len(log.Topics))
	}

	fields := make(map[string]interface{})
	err := abis.ERC1155.UnpackIntoMap(fields, erc1155Batch, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid number of fields (want: %d, have: %d)", 2, len(fields))
	}

	fieldIDs, ok := fields[erc1155IDs]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", erc1155IDs)
	}

	fieldValues, ok := fields["values"]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", erc1155Values)
	}

	tokenIDs, ok := fieldIDs.([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type (field: %s, have: %T)", erc1155IDs, fieldIDs)
	}

	values, ok := fieldValues.([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type (field: %s, have: %T)", erc1155Values, fieldValues)
	}

	if len(tokenIDs) != len(values) {
		return nil, fmt.Errorf("length mismatch between ids and values fields (ids: %d, values: %d)", len(tokenIDs), len(values))
	}

	var transfers []*events.Transfer
	for index, tokenID := range tokenIDs {
		value := values[index]

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
			TokenCount:        value.String(),
			// EmittedAt set after parsing
		}
		transfers = append(transfers, &transfer)
	}

	return transfers, nil
}
