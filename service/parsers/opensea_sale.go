package parsers

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"

	"github.com/NFT-com/indexer/models/events"
)

func OpenSeaSale(log types.Log) (*events.Sale, error) {

	fields := make(map[string]interface{})
	err := abis.OpenSea.UnpackIntoMap(fields, "OrdersMatched", log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	price, ok := fields["price"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for \"id\" field (%T)", fields["price"])
	}

	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
	hash := sha3.Sum256(data)

	sale := events.Sale{
		ID:                 hex.EncodeToString(hash[:]),
		MarketplaceAddress: log.Address.Hex(),
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TxHash.Hex(),
		EventIndex:         log.Index,
		SellerAddress:      log.Topics[1].Hex(),
		BuyerAddress:       log.Topics[2].Hex(),
		TradePrice:         price.String(),
		// EmmittedAt set after processing
	}

	return &sale, nil
}