package parsers

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"github.com/NFT-com/indexer/models/hashes"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/NFT-com/indexer/models/events"
)

const (
	defaultEthereumCurrency = "ETH"
	openSeaAddress          = "0x7f268357a8c2552623316e2562d90e642bb538e5"
)

func OpenSeaSale(logs []types.Log) (*events.Sale, error) {

	log := openSeaLog(logs)
	if log == nil {
		return nil, fmt.Errorf("could not find opensea log")
	}

	fields := make(map[string]interface{})
	err := abis.OpenSea.UnpackIntoMap(fields, "OrdersMatched", log.Data)
	if err != nil {
		return nil, fmt.Errorf("could not unpack log fields: %w", err)
	}

	price, ok := fields["price"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type for \"price\" field (%T)", fields["price"])
	}

	data := make([]byte, 8+32+8)
	binary.BigEndian.PutUint64(data[0:8], log.BlockNumber)
	copy(data[8:40], log.TxHash[:])
	binary.BigEndian.PutUint64(data[40:48], uint64(log.Index))
	hash := sha3.Sum256(data)
	saleID := uuid.Must(uuid.FromBytes(hash[:16]))

	collectionAddress, tokenID, err := retrieveTokenInformation(logs)
	if err != nil {
		return nil, fmt.Errorf("could retrieve token information: %w", err)
	}

	tradeCurrency, err := retrieveTradeCurrency(logs, price)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve trade currency: %w", err)
	}

	sale := events.Sale{
		ID: saleID.String(),
		// ChainID set after parsing
		MarketplaceAddress: log.Address.Hex(),
		CollectionAddress:  collectionAddress,
		TokenID:            tokenID,
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TxHash.Hex(),
		EventIndex:         log.Index,
		SellerAddress:      log.Topics[1].Hex(),
		BuyerAddress:       log.Topics[2].Hex(),
		TradePrice:         price.String(),
		TradeCurrency:      tradeCurrency,
		// EmmittedAt set after parsing
	}

	return &sale, nil
}

func openSeaLog(logs []types.Log) *types.Log {
	for _, log := range logs {

		if strings.ToLower(log.Address.Hex()) == strings.ToLower(openSeaAddress) {
			return &log
		}
	}

	return nil
}

func retrieveTokenInformation(logs []types.Log) (string, string, error) {
	for _, log := range logs {

		// for now, we only care about single transactions as I didn't see any batch transfer being used
		switch log.Topics[0].String() {

		case hashes.ERC721Transfer:

			event, err := ERC721Transfer(log)
			if err != nil {
				return "", "", fmt.Errorf("could not parse event: %w", err)
			}

			// check if the event has the price value set as it can be an erc20 transfer (same hash)
			if event.Value != "" {
				continue
			}

			return event.CollectionAddress, event.TokenID, nil

		case hashes.ERC1155Transfer:

			event, err := ERC1155Transfer(log)
			if err != nil {
				return "", "", fmt.Errorf("could not parse event: %w", err)
			}

			// check if the event has the price value set as it can be an erc20 transfer (same hash)
			if event.Value != "" {
				continue
			}

			return event.CollectionAddress, event.TokenID, nil
		}

	}

	return "", "", fmt.Errorf("token information not found")
}

func retrieveTradeCurrency(logs []types.Log, price *big.Int) (string, error) {
	for _, log := range logs {

		if log.Topics[0].String() != hashes.ERC20Transfer {
			continue
		}

		event, err := ERC20Transfer(log)
		if err != nil {
			return "", fmt.Errorf("could not parse event: %w", err)
		}

		// check if the event has the price value set and equal to price as it can be an erc721 transfer (same hash)
		if event.Value != price.String() {
			continue
		}

		return event.CollectionAddress, nil
	}

	return defaultEthereumCurrency, nil
}
