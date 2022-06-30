package id

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

func Event(hash string, index uint) string {
	eventHash := sha3.Sum256([]byte(fmt.Sprintf("%s-%d", hash, index)))
	eventID := uuid.Must(uuid.FromBytes(eventHash[:]))
	return eventID.String()
}
