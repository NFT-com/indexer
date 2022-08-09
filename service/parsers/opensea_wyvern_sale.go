package parsers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
)

const (
	eventOrdersMatched = "OrdersMatched"
	fieldPrice         = "price"
)

func OpenSeaWyvernSale(log types.Log) (*events.Sale, error) {

	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("invalid topic lenght have (%d) want (%d)", len(log.Topics), 3)
	}

	fields := make(map[string]interface{})
	err := abis.OpenSeaWyvern.UnpackIntoMap(fields, eventOrdersMatched, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	price, ok := fields[fieldPrice].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for %q field (%T)", fieldPrice, fields[fieldPrice])
	}

	sale := events.Sale{
		ID: logID(log),
		// ChainID set after parsing
		MarketplaceAddress: log.Address.Hex(),
		CollectionAddress:  "", // Done in completion pipeline
		TokenID:            "", // Done in completion pipeline
		TokenCount:         0,  // Done in completion pipeline
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
