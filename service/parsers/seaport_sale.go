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
	return n.Address != addressZero && n.Identifier != nil && n.Amount != nil
}

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

	offerer := common.BytesToAddress(log.Topics[1].Bytes())

	_, ok := fields[seaportOrder]
	if !ok {
		return nil, fmt.Errorf("missing field key (%s)", seaportOrder)
	}

	_, hasFulfiller := fields[seaportFulfiller]
	_, hasRecipient := fields[seaportRecipient]
	if !hasFulfiller && !hasRecipient {
		return nil, fmt.Errorf("missing field key (%s or %s)", seaportFulfiller, seaportRecipient)
	}

	fieldOffer, ok := fields[seaportOffer]
	if !ok {
		return nil, fmt.Errorf("missing field key (%s)", seaportOffer)
	}

	fieldConsideration, ok := fields[seaportConsideration]
	if !ok {
		return nil, fmt.Errorf("missing field key (%s)", seaportConsideration)
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

	// We are looking for 3 required and 1 optional component:
	// - NFT;
	// - payment;
	// - fee;
	// - tip.
	var nft NFT
	var payment Transfer
	var fee Transfer
	var tip Transfer

	// The offer corresponds of either to the NFT or the payment.
	offerItem := offerItems[0]
	switch offerItem.ItemType {

	// 0 - native token
	// 1 - ERC20 token
	case 0, 1:
		payment = Transfer{
			Address: offerItem.Token,
			Amount:  offerItem.Amount,
		}

	// 2 - ERC721 token
	// 3 - ERC1155 token
	case 2, 3:
		nft = NFT{
			Address:    offerItem.Token,
			Identifier: offerItem.Identifier,
			Amount:     offerItem.Amount,
		}

	// 4 - ERC721 token with extra criteria
	// 5 - ERC1155 token with extra criteria
	case 4, 5:
		return nil, fmt.Errorf("unsupported sale (additional criteria in offer)")

	default:
		return nil, fmt.Errorf("unsupported sale (unknown offer item type %d)", offerItem.ItemType)
	}

	// After identifying the offered NFT, we categorize the considerations.
	for _, item := range considerationItems {

		// If we have two NFTs, we don't support this sale.
		if (item.ItemType == 2 || item.ItemType == 3) && nft.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple non-fungible tokens)")
		}

		// If we have two payments, we don't support this sale.
		if (item.ItemType == 0 || item.ItemType == 1) && item.Recipient == offerer && payment.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple payments)")
		}

		// If we have to fees, we don't support this sale.
		if (item.ItemType == 0 || item.ItemType == 1) && item.Recipient == addressFee && fee.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple fees)")
		}

		// If we have multiple tips, we don't support this sale.
		if (item.ItemType == 0 || item.ItemType == 1) && item.Recipient != offerer && item.Recipient != addressFee && tip.Valid() {
			return nil, fmt.Errorf("unsupported sale (multiple tips)")
		}

		switch item.ItemType {

		// 0 - native token
		// 1 - ERC20 token
		case 0, 1:

			transfer := Transfer{
				Address: item.Token,
				Amount:  item.Amount,
			}

			switch item.Recipient {
			case offerer:
				payment = transfer
			case addressFee:
				fee = transfer
			default:
				tip = transfer
			}

		// 2 - ERC721 token
		// 3 - ERC1155 token
		case 2, 3:
			nft = NFT{
				Address:    item.Token,
				Identifier: item.Identifier,
				Amount:     item.Amount,
			}

		// 4 - ERC721 token with extra criteria
		// 5 - ERC1155 token with extra criteria
		case 4, 5:
			return nil, fmt.Errorf("unsupported sale (additional criteria in consideration)")

		default:
			return nil, fmt.Errorf("unsupported sale (unknown consideration item type %d)", offerItem.ItemType)
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
