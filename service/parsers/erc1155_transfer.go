package parsers

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
)

func ERC1155Transfer(log types.Log) (*events.Transfer, error) {

	fields := make(map[string]interface{})
	err := abis.ERC1155.UnpackIntoMap(fields, "TransferSingle", log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	tokenID, ok := fields["id"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for \"id\" field (%T)", fields["id"])
	}
	count, ok := fields["value"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for \"value\" field (%T)", fields["value"])
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
		CollectionAddress: log.Address.Hex(),
		TokenID:           tokenID.String(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		SenderAddress:     log.Topics[2].Hex(),
		ReceiverAddress:   log.Topics[3].Hex(),
		TokenCount:        count.Uint64(),
		// EmmittedAt set after parsing
	}

	return &transfer, nil
}
