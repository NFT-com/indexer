package parsers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
)

type offer struct {
	ItemType   uint8          `json:"itemType"`
	Token      common.Address `json:"token"`
	Identifier *big.Int       `json:"identifier"`
	Amount     *big.Int       `json:"amount"`
}

type consideration struct {
	ItemType   uint8          `json:"itemType"`
	Token      common.Address `json:"token"`
	Identifier *big.Int       `json:"identifier"`
	Amount     *big.Int       `json:"amount"`
	Recipient  common.Address `json:"recipient"`
}

func OpenSeaSeaportSale(log types.Log) (*events.Sale, error) {

	fields := make(map[string]interface{})
	err := abis.OpenSeaSeaport.UnpackIntoMap(fields, "OrderFulfilled", log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	// Get the buys from the event.
	recipient, ok := fields["recipient"].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid type for \"recipient\" field (%T)", fields["recipient"])
	}

	// Retrieve the offers from the event.
	offers := make([]offer, 0)
	err = getCompositeData(fields["offer"], &offers)
	if err != nil {
		return nil, fmt.Errorf("could not get \"offer\" field: %w", err)
	}

	// Currently we will ignore all events with multiple currencies.
	if len(offers) > 1 {
		return nil, fmt.Errorf("could not parse event: multiple offers not supported")
	}

	offer := offers[0]

	// Retrieve consideration items.
	considerations := make([]consideration, 0)
	err = getCompositeData(fields["consideration"], &considerations)
	if err != nil {
		return nil, fmt.Errorf("could not get \"consideration\" field: %w", err)
	}

	if len(considerations) == 0 {
		return nil, fmt.Errorf("could not get considerations: considerations are empty")
	}

	// filter out fees paid to the opensea market
	considerations = filterFees(considerations, offer.Token, offer.Identifier)

	if isSaleOrder(offer) {
		considerations = append(considerations[:1], filterFees(considerations[1:], considerations[0].Token, considerations[0].Identifier)...)
	}

	// Currently we will ignore all events with multiple tokens sold.
	if len(considerations) > 1 {
		return nil, fmt.Errorf("could not parse event: multiple considerations per sale not supported")
	}

	consideration := considerations[0]

	switch {

	// in this case the offer var represents the nft being sold and the consideration represents the payment for it.
	case isSaleOrder(offer):

		sale := events.Sale{
			ID: saleID(log),
			// ChainID set after parsing
			MarketplaceAddress: log.Address.Hex(),
			CollectionAddress:  offer.Token.Hex(),
			TokenID:            offer.Identifier.String(),
			TokenCount:         uint(offer.Amount.Uint64()),
			BlockNumber:        log.BlockNumber,
			TransactionHash:    log.TxHash.Hex(),
			EventIndex:         log.Index,
			SellerAddress:      common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
			BuyerAddress:       recipient.Hex(),
			CurrencyAddress:    consideration.Token.String(),
			CurrencyValue:      consideration.Amount.String(),
			// EmittedAt set after parsing
			NeedsCompletion: false,
		}

		return &sale, nil

	default:
		// in this case the consideration var represents the nft being sold and the offer the payment for it.

		sale := events.Sale{
			ID: saleID(log),
			// ChainID set after parsing
			MarketplaceAddress: log.Address.Hex(),
			CollectionAddress:  consideration.Token.Hex(),
			TokenID:            consideration.Identifier.String(),
			TokenCount:         uint(consideration.Amount.Uint64()),
			BlockNumber:        log.BlockNumber,
			TransactionHash:    log.TxHash.Hex(),
			EventIndex:         log.Index,
			SellerAddress:      common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
			BuyerAddress:       recipient.Hex(),
			CurrencyAddress:    offer.Token.String(),
			CurrencyValue:      offer.Amount.String(),
			// EmittedAt set after parsing
			NeedsCompletion: false,
		}

		return &sale, nil
	}
}

func isSaleOrder(offer offer) bool {
	// ItemTypes:
	// 0: ETH on mainnet, MATIC on polygon, etc.
	// 1: ERC20 items (ERC777 and ERC20 analogues could also technically work)
	// 2: ERC721 items
	// 3: ERC1155 items
	// 4: ERC721 items where a number of tokenIds are supported
	// 5: ERC1155 items where a number of ids are supported
	// if the offer type is greater than two means this is a sell order not a buy order, meaning the offer will be the
	// nft instead of the payment
	return offer.ItemType >= 2
}

func getCompositeData(field interface{}, out interface{}) error {
	data, err := json.Marshal(field)
	if err != nil {
		return fmt.Errorf("could not marshal composite data: %w", err)
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("could not unmarshal composite data: %w", err)
	}

	return nil
}

func filterFees(considerations []consideration, token common.Address, identifier *big.Int) []consideration {
	filtered := make([]consideration, 0)

	for _, consideration := range considerations {
		// if the contract addresses are the same remove it as it represents fees
		if consideration.Token.Hex() == token.Hex() &&
			consideration.Identifier.Cmp(identifier) == 0 {
			continue
		}

		filtered = append(filtered, consideration)
	}

	return filtered
}

func saleID(log types.Log) string {
	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
	hash := sha3.Sum256(data)
	saleID := uuid.Must(uuid.FromBytes(hash[:16]))

	return saleID.String()
}
