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
	seaportMatch         = "OrderFulfilled"
	seaportOrder         = "orderHash"
	seaportRecipient     = "recipient"
	seaportFulfiller     = "fulfiller"
	seaportOffer         = "offer"
	seaportConsideration = "consideration"

	itemType       = "itemType"
	itemToken      = "token"
	itemIdentifier = "identifier"
	itemAmount     = "amount"
	itemRecipient  = "recipient"
)

type Token struct {
	Address    common.Address
	Identifier *big.Int
	Amount     *big.Int
}

func SeaportSale(log types.Log) (*events.Sale, error) {

	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("invalid number of topics (want: %d, have: %d)", 3, len(log.Topics))
	}

	fields := make(map[string]interface{})
	err := abis.OpenSeaSeaport.UnpackIntoMap(fields, seaportMatch, log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	if len(fields) != 4 {
		return nil, fmt.Errorf("invalid number of fields (want: %d, have: %d)", 4, len(fields))
	}

	_, ok := fields[seaportOrder]
	if !ok {
		return nil, fmt.Errorf("missing field key (%s)", seaportOrder)
	}

	fieldOfferer, ok := fields[seaportRecipient]
	if !ok {
		fieldOfferer, ok = fields[seaportFulfiller]
		if !ok {
			return nil, fmt.Errorf("missing field key (%s or %s)", seaportRecipient, seaportFulfiller)
		}
	}

	fieldOffer, ok := fields[seaportOffer]
	if !ok {
		return nil, fmt.Errorf("missing field key (%s)", seaportOffer)
	}

	fieldConsideration, ok := fields[seaportConsideration]
	if !ok {
		return nil, fmt.Errorf("missing field key (%s)", seaportConsideration)
	}

	offerer, ok := fieldOfferer.(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid field type (field: %s, want: %T, have: %T)", seaportRecipient, common.Address{}, fieldOfferer)
	}

	offerItems, ok := fieldOffer.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid field type (field: %s, want: %T, have: %T)", seaportOffer, []interface{}{}, fieldOffer)
	}

	considerationItems, ok := fieldConsideration.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid field type (field: %s, want: %T, have: %T)", seaportConsideration, []interface{}{}, fieldConsideration)
	}

	var nfts []Token
	var fts []Token

	for _, item := range offerItems {

		lookup, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid offer item type (want: %T, have: %T)", map[string]interface{}{}, item)
		}

		fieldType, ok := lookup[itemType]
		if !ok {
			return nil, fmt.Errorf("missing offer field key (%s)", itemType)
		}

		fieldContract, ok := lookup[itemToken]
		if !ok {
			return nil, fmt.Errorf("missing offer field key (%s)", itemToken)
		}

		fieldIdentifier, ok := lookup[itemIdentifier]
		if !ok {
			return nil, fmt.Errorf("missing offer field key (%s)", itemIdentifier)
		}

		fieldAmount, ok := lookup[itemAmount]
		if !ok {
			return nil, fmt.Errorf("missing offer field key (%s)", itemAmount)
		}

		typ, ok := fieldType.(uint8)
		if !ok {
			return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemType, uint8(0), fieldType)
		}

		contract, ok := fieldContract.(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemToken, common.Address{}, fieldContract)
		}

		identifier, ok := fieldIdentifier.(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemIdentifier, &big.Int{}, fieldIdentifier)
		}

		amount, ok := fieldAmount.(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemAmount, &big.Int{}, fieldAmount)
		}

		token := Token{
			Address:    contract,
			Identifier: identifier,
			Amount:     amount,
		}

		switch typ {
		case 0, 1:
			fts = append(fts, token)
		case 2, 3, 4, 5:
			nfts = append(nfts, token)
		}
	}

	if len(nfts) == 0 && len(fts) == 0 {
		return nil, fmt.Errorf("unsupported event (offer contains no items)")
	}

	if len(nfts) > 0 && len(fts) > 0 {
		return nil, fmt.Errorf("unsupported event (offer contains fungible and non-fungible token items)")
	}

	for _, item := range considerationItems {

		lookup, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid consideration item type (want: %T, have: %T)", map[string]interface{}{}, item)
		}

		fieldType, ok := lookup[itemType]
		if !ok {
			return nil, fmt.Errorf("missing consideration field key (%s)", itemType)
		}

		fieldToken, ok := lookup[itemToken]
		if !ok {
			return nil, fmt.Errorf("missing consideration field key (%s)", itemToken)
		}

		fieldIdentifier, ok := lookup[itemIdentifier]
		if !ok {
			return nil, fmt.Errorf("missing consideration field key (%s)", itemIdentifier)
		}

		fieldRecipient, ok := lookup[itemRecipient]
		if !ok {
			return nil, fmt.Errorf("missing consideration field key (%s)", itemRecipient)
		}

		fieldAmount, ok := lookup[itemAmount]
		if !ok {
			return nil, fmt.Errorf("missing consideration field key (%s)", itemAmount)
		}

		typ, ok := fieldType.(uint8)
		if !ok {
			return nil, fmt.Errorf("invalid consideration field type (field: %s, want: %T, have: %T)", itemType, uint8(0), fieldType)
		}

		token, ok := fieldToken.(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid consideration field type (field: %s, want: %T, have: %T)", itemToken, common.Address{}, fieldToken)
		}

		identifier, ok := fieldIdentifier.(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemIdentifier, &big.Int{}, fieldIdentifier)
		}

		amount, ok := fieldAmount.(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid consideration field type (field: %s, want: %T, have: %T)", itemAmount, &big.Int{}, fieldAmount)
		}

		recipient, ok := fieldRecipient.(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid consideration field type (field: %s, want: %T, have: %T)", itemRecipient, common.Address{}, fieldRecipient)
		}

		if recipient != offerer {
			continue
		}

		item := Token{
			Address:    token,
			Identifier: identifier,
			Amount:     amount,
		}

		switch typ {
		case 0, 1:
			fts = append(fts, item)
		case 2, 3, 4, 5:
			nfts = append(nfts, item)
		}
	}

	if len(fts) == 0 {
		return nil, fmt.Errorf("unsupported event (no fungible token transfers")
	}

	if len(nfts) == 0 {
		return nil, fmt.Errorf("unsupported event (no non-fungible token transfers")
	}

	if len(fts) > 1 {
		return nil, fmt.Errorf("unsupported event (multiple fungible token transfers)")
	}

	if len(nfts) > 1 {
		return nil, fmt.Errorf("unsupported event (multiple non-fungible token transfers)")
	}

	ft := fts[0]
	nft := nfts[0]

	sale := events.Sale{
		ID: id.Log(log),
		// ChainID set after parsing
		MarketplaceAddress: log.Address.Hex(),
		CollectionAddress:  nft.Address.Hex(),
		TokenID:            nft.Identifier.String(),
		TokenCount:         uint(nft.Amount.Uint64()),
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TxHash.Hex(),
		EventIndex:         log.Index,
		SellerAddress:      common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
		BuyerAddress:       offerer.Hex(),
		CurrencyAddress:    ft.Address.String(),
		CurrencyValue:      ft.Amount.String(),
		// EmittedAt set after parsing
		NeedsCompletion: false,
	}

	return &sale, nil
}
