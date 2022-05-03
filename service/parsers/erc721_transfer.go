package parsers

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"

	"github.com/NFT-com/indexer/models/events"
)

func ERC721Transfer(log types.Log) (*events.Transfer, error) {

	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
	hash := sha3.Sum256(data)

	transfer := events.Transfer{
		ID:                hex.EncodeToString(hash[:]),
		CollectionAddress: log.Address.Hex(),
		BaseTokenID:       "",
		TokenID:           log.Topics[3].Big().String(),
		BlockNumber:       log.BlockNumber,
		EventIndex:        log.Index,
		TransactionHash:   log.TxHash.Hex(),
		FromAddress:       log.Topics[1].Hex(),
		ToAddress:         log.Topics[2].Hex(),
		// EmmittedAt set after processing
	}

	return &transfer, nil
}
