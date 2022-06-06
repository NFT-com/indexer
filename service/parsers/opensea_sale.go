package parsers

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
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
	saleID := uuid.Must(uuid.FromBytes(hash[:16]))

	sale := events.Sale{
		ID: saleID.String(),
		// ChainID set after parsing
		MarketplaceAddress: log.Address.Hex(),
		CollectionAddress:  "", // TODO
		TokenID:            "", // TODO
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TxHash.Hex(),
		EventIndex:         log.Index,
		SellerAddress:      common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
		BuyerAddress:       common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		TradePrice:         price.String(),
		// EmittedAt set after parsing
	}

	return &sale, nil
}
