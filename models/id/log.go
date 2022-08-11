package id

import (
	"encoding/binary"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/core/types"
)

func Log(log types.Log) string {

	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))

	hash := sha3.Sum256(data)
	id := uuid.Must(uuid.FromBytes(hash[:16]))

	return id.String()
}
