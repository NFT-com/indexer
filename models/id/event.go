package id

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

func Event(tx string, index uint) string {
	eventHash := sha3.Sum256([]byte(fmt.Sprintf("%s-%d", tx, index)))
	eventID := uuid.Must(uuid.FromBytes(eventHash[:16]))
	return eventID.String()
}
