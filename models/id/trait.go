package id

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

func Trait(chainID uint64, address string, tokenID string, index uint) string {
	traitHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s-%d", chainID, address, tokenID, index)))
	traitID := uuid.Must(uuid.FromBytes(traitHash[:16]))
	return traitID.String()
}
