package id

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

func NFT(chainID uint64, address string, tokenID string) string {
	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s", chainID, address, tokenID)))
	nftID := uuid.Must(uuid.FromBytes(nftHash[:16]))
	return nftID.String()
}
