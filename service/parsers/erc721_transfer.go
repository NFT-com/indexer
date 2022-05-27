package parsers

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/NFT-com/indexer/models/events"
)

func ERC721Transfer(log types.Log) (*events.Transfer, error) {

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
		TokenID:           log.Topics[3].Big().String(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		SenderAddress:     log.Topics[1].Hex(),
		ReceiverAddress:   log.Topics[2].Hex(),
		TokenCount:        1,
		// EmittedAt set after parsing
	}

	return &transfer, nil
}
