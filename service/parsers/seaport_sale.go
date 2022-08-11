package parsers

import (
	"encoding/json"
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
)

var (
	addressFee = common.HexToAddress("0x8De9C5A032463C561423387a9648c5C7BCC5BC90")
)

type Transfer struct {
	Address common.Address
	Amount  *big.Int
}

func (t Transfer) Valid() bool {
	return t.Amount != nil
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

	out, _ := json.Marshal(fields)
	fmt.Println(string(out))

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
		return nil, fmt.Errorf("invalid field type (field: %s, have: %T)", seaportRecipient, fieldOfferer)
	}

	// This is a bit messy, but unfortunately, we can't type assert slices with concrete named types when the decoded
	// type is a slice of anonymous structs.
	offerItems, ok := fieldOffer.([]struct {
		ItemType   uint8          `json:"itemType"`
		Token      common.Address `json:"token"`
		Identifier *big.Int       `json:"identifier"`
		Amount     *big.Int       `json:"amount"`
	})
	if !ok {
		return nil, fmt.Errorf("invalid field type (field: %s, have: %T)", seaportOffer, fieldOffer)
	}
	considerationItems, ok := fieldConsideration.([]struct {
		ItemType   uint8          `json:"itemType"`
		Token      common.Address `json:"token"`
		Identifier *big.Int       `json:"identifier"`
		Amount     *big.Int       `json:"amount"`
		Recipient  common.Address `json:"recipient"`
	})
	if !ok {
		return nil, fmt.Errorf("invalid field type (field: %s, have: %T)", seaportConsideration, fieldConsideration)
	}

	// In general, basic offers should only have a single item in their offer.
	if len(offerItems) == 0 {
		return nil, fmt.Errorf("invalid sale (no offer items)")
	}
	if len(offerItems) > 1 {
		return nil, fmt.Errorf("unsupported sale (multiple offer items")
	}

	// 2 and 3 correspond to ERC721 and ERC1155 NFT tokens respectively.
	// We only support offers of NFTs without extra conditions, so that's what we want.
	nft := offerItems[0]

	// 0 and 1 correspond to native and ERC20 fungible tokens respectively.
	// We don't support fungible tokens in the offer.
	if nft.ItemType == 0 || nft.ItemType == 1 {
		return nil, fmt.Errorf("unsupported sale (fungible token in offer)")
	}

	// 4 and 5 correspond to ERC721 and ERC1155 NFT tokens with additional sale criteria.
	// We don't support additional criteria, as it makes the price untractable.
	if nft.ItemType == 4 || nft.ItemType == 5 {
		return nil, fmt.Errorf("unsupported sale (additional criteria in offer)")
	}

	// After identifying the offered NFT, we categorize the considerations.
	var payment Transfer
	var fee Transfer
	var tip Transfer
	for _, item := range considerationItems {

		// 0 and 1 correspond to native and ERC20 fungible tokens respectively.
		// That's what we need every consideration to be.
		if item.ItemType != 0 && item.ItemType != 1 {
			return nil, fmt.Errorf("unsupported sale (non-fungible token in consideration)")
		}

		transfer := Transfer{
			Address: item.Token,
			Amount:  item.Amount,
		}

		switch {

		// If the recipient of a consideration is the offerer, this is the payment for the NFT.
		case item.Recipient == offerer:
			payment = transfer

		// If the recipient of a consideration is the OpenSea fee address, it's the fee for the sale.
		case item.Recipient == addressFee:
			fee = transfer

		// Otherwise, it's also possible that there is exactly on optional tip.
		default:
			tip = transfer
		}
	}

	// We need at least a valid NFT, a valid payment and a valid fee.
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
		CollectionAddress:  nft.Token.Hex(),
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
