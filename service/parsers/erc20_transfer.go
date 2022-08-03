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
	eventTransfer = "Transfer"
	fieldValue    = "value"
)

func ERC20Transfer(log types.Log) (*events.Transfer, error) {

	fields := make(map[string]interface{})
	err := abis.ERC20.UnpackIntoMap(fields, eventTransfer, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	value, ok := fields[fieldValue].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for %q field (%T)", fieldValue, fields[fieldValue])
	}

	transfer := events.Transfer{
		ID: id(log),
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

func id(log types.Log) string {
	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
	hash := sha3.Sum256(data)
	transferID := uuid.Must(uuid.FromBytes(hash[:16]))

	return transferID.String()
}
