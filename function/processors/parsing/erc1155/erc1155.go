package erc1155

import (
	"math/big"
)

// ExtractIDs takes a `composed ID` and returns the split collection and NFT IDs.
func ExtractIDs(composedID []byte) (*big.Int, *big.Int) {
	// If the colID is 0 then it only had a NFT ID
	if len(composedID) <= 16 {
		nftID := big.NewInt(0)
		nftID.SetBytes(composedID[:])
		return nil, nftID
	}

	// Big Int is composed of two uint64
	firstUint64LastByte := len(composedID) - 16

	colID := big.NewInt(0)
	colID.SetBytes(composedID[:firstUint64LastByte])

	nftID := big.NewInt(0)
	nftID.SetBytes(composedID[firstUint64LastByte:])

	return colID, nftID
}
