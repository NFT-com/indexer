package parsers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/id"
)

const (
	wyvernMatch = "OrdersMatched"
	wyvernBuy   = "buyHash"
	wyvernSell  = "sellHash"
	wyvernPrice = "price"
)

func WyvernSale(log types.Log) (*events.Sale, error) {

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid number of topics (want: %d, have: %d)", 4, len(log.Topics))
	}

	fields := make(map[string]interface{})
	err := abis.OpenSeaWyvern.UnpackIntoMap(fields, wyvernMatch, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	if len(fields) != 3 {
		return nil, fmt.Errorf("invalid number of fields (want: %d, have: %d)", 3, len(fields))
	}

	_, ok := fields[wyvernBuy]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", wyvernBuy)
	}

	_, ok = fields[wyvernSell]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", wyvernSell)
	}

	fieldPrice, ok := fields[wyvernPrice]
	if !ok {
		return nil, fmt.Errorf("missing field (%s)", wyvernPrice)
	}

	price, ok := fieldPrice.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type (field: %s,  have: %T)", wyvernPrice, fieldPrice)
	}

	sale := events.Sale{
		ID: id.Log(log),
		// ChainID set after parsing
		MarketplaceAddress: log.Address.Hex(),
		CollectionAddress:  "",  // Done in completion pipeline
		TokenID:            "",  // Done in completion pipeline
		TokenCount:         "0", // Done in completion pipeline
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TxHash.Hex(),
		EventIndex:         log.Index,
		SellerAddress:      common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
		BuyerAddress:       common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		CurrencyAddress:    "", // Done in completion pipeline
		CurrencyValue:      price.String(),
		// EmittedAt set after parsing
		NeedsCompletion: true,
	}

	return &sale, nil
}
