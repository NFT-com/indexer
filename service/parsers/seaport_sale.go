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

var (
	addressZero = common.Address{}
	addressFee  = common.HexToAddress("0x8De9C5A032463C561423387a9648c5C7BCC5BC90")
)

type NFT struct {
	Address    common.Address
	Identifier *big.Int
	Amount     *big.Int
}

func (n NFT) Valid() bool {
	return n.Address != addressZero && n.Amount.Cmp(big.NewInt(0)) != 0
}

type Transfer struct {
	Address common.Address
	Amount  *big.Int
}

func (t Transfer) Valid() bool {
	return t.Amount.Cmp(big.NewInt(0)) != 0
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

	// A simple order on OpenSea has 3 or 4 components:
	// - the NFT that is sold;
	// - the payment in ERC20 or native token;
	// - a fee payment in the same token; and
	// - an optional tip in the same token.
	// We try to identify each of the four components; if anything is off, the offer
	// structure is too complex and we don't handle the edge case for now.
	var nft NFT
	var payment Transfer
	var tip Transfer
	var fee Transfer

	// In general, basic offers should only have a single item in their offer.
	if len(offerItems) == 0 {
		return nil, fmt.Errorf("invalid sale (no offer items)")
	}
	if len(offerItems) > 1 {
		return nil, fmt.Errorf("unsupported sale (multiple offer items")
	}

	// Next, we identify whether the offerer is putting up an NFT for sale, or if he
	// is offering a payment for a certain NFT.
	item := offerItems[0]
	lookup, ok := item.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid offer item type (want: %T, have: %T)", map[string]interface{}{}, item)
	}

	fieldType, ok := lookup[itemType]
	if !ok {
		return nil, fmt.Errorf("missing offer field key (%s)", itemType)
	}

	fieldToken, ok := lookup[itemToken]
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

	token, ok := fieldToken.(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemToken, common.Address{}, fieldToken)
	}

	identifier, ok := fieldIdentifier.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemIdentifier, &big.Int{}, fieldIdentifier)
	}

	amount, ok := fieldAmount.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid order field type (field: %s, want: %T, have: %T)", itemAmount, &big.Int{}, fieldAmount)
	}

	// 0 and 1 correspond to native and ERC20 fungible tokens respectively.
	// 2 and 3 correspond to ERC721 and ERC1155 NFT tokens respectively.
	// 4 and 5 correspond to ERC721 and ERC1155 NFT tokens with additional sale criteria.
	switch typ {

	case 0, 1:
		payment = Transfer{
			Address: token,
			Amount:  amount,
		}

	case 2, 3:
		nft = NFT{
			Address:    token,
			Identifier: identifier,
			Amount:     amount,
		}

	case 4, 5:
		return nil, fmt.Errorf("unsupported sale (additional offer criteria)")

	default:
		return nil, fmt.Errorf("unknown item type (%d)", typ)
	}

	// After identifying the offer, we look at the considerations to classify them.
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

		// We don't support extra criteria for now.
		if typ == 4 || typ == 5 {
			return nil, fmt.Errorf("unsupported sale (additional consideration criteria)")
		}

		// If there is a NFT, but we already had it, the trade has multiple NFTs and is unsupported.
		if (typ == 2 || typ == 3) && nft.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple NFTs)")
		}

		// If there is a payment, but we already have one, the trade has multiple currencies and is unsupported.
		if (typ == 0 || typ == 1) && recipient == offerer && payment.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple payments)")
		}

		// If there is a fee, but we already have one, the trade has multiple fees and is unsupported.
		if (typ == 0 || typ == 1) && recipient == addressFee && fee.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple fees)")
		}

		// If there is a tip, but we already have one, the trade has multiple tips and is unsupported.
		if (typ == 0 || typ == 1) && recipient != offerer && recipient != addressFee && tip.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple fees)")
		}

		// At this point, we can extract the component depending on conditions.
		switch {

		case typ == 2 || typ == 3:
			nft = NFT{
				Address:    token,
				Identifier: identifier,
				Amount:     amount,
			}

		case recipient == offerer:
			payment = Transfer{
				Address: token,
				Amount:  amount,
			}

		case recipient == addressFee:
			fee = Transfer{
				Address: token,
				Amount:  amount,
			}

		default:
			tip = Transfer{
				Address: token,
				Amount:  amount,
			}
		}
	}

	// We need at least a valid NFT, a valid payment and a valid fee.
	if !nft.Valid() {
		return nil, fmt.Errorf("unsupported sale (no NFT)")
	}
	if !payment.Valid() {
		return nil, fmt.Errorf("unsupported sale (no payment)")
	}
	if !fee.Valid() {
		return nil, fmt.Errorf("unsupported sale (no fee)")
	}

	// We need the token types of all token transfers to match.
	if fee.Address != payment.Address {
		return nil, fmt.Errorf("unsupported sale (fee and payment token mismatch)")
	}
	if tip.Valid() && tip.Address != payment.Address {
		return nil, fmt.Errorf("unsupported sale (tip and payment token mismatch)")
	}

	// The fee and tip should also be smaller than the payment.
	if fee.Amount.Cmp(payment.Amount) > 0 {
		return nil, fmt.Errorf("invalid sale (fee bigger than payment)")
	}
	if tip.Valid() && tip.Amount.Cmp(payment.Amount) > 0 {
		return nil, fmt.Errorf("invalid sale (tip bigger than payment)")
	}

	// At this point, we know what the NFT is and what the payment is. Tip and fee
	// can be ignored with the current data model.
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
		CurrencyAddress:    payment.Address.String(),
		CurrencyValue:      payment.Amount.String(),
		// EmittedAt set after parsing
		NeedsCompletion: false,
	}

	return &sale, nil
}
