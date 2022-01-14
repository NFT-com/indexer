package contracts

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/store"
)

const (
	TopicTransfer      = "Transfer"
	TopicOrdersMatched = "OrdersMatched"
)

type Contract struct {
	log zerolog.Logger

	client        *ethclient.Client
	contractStore store.ContractStore
}

func New(log zerolog.Logger, client *ethclient.Client, contractStore store.ContractStore) *Contract {
	c := Contract{
		log:           log.With().Str("component", "contract_parser").Logger(),
		client:        client,
		contractStore: contractStore,
	}

	return &c
}

func (c *Contract) ParseEvent(ctx context.Context, log types.Log) (events.Event, error) {
	topic := log.Topics[0]

	stringAbi, err := c.contractStore.GetABI(ctx, "", "", log.Address.Hex())
	if err != nil {
		return nil, err
	}

	parsedAbi, err := abi.JSON(strings.NewReader(stringAbi))
	if err != nil {
		return nil, err
	}

	event, err := parsedAbi.EventByID(topic)
	if err != nil {
		return nil, err
	}

	switch event.Name {
	case TopicTransfer:
		return c.parseTransferEvent(parsedAbi, log)
	case TopicOrdersMatched:
		return c.parseOrdersMatchedEvent(parsedAbi, log)
	default:
		// FIXME: Handle more topics.
	}

	return nil, nil
}

func (c *Contract) parseTransferEvent(abi abi.ABI, log types.Log) (events.Event, error) {
	data, err := abi.Unpack(TopicTransfer, log.Data)
	if err != nil {
		c.log.Error().Err(err).Str("topic", TopicTransfer).Msg("Failed to unpack event")
		return nil, err
	}

	from, ok := data[0].(common.Address)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	to, ok := data[1].(common.Address)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	nftID, ok := data[2].(*big.Int)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	return events.NewTransfer(
		log.TxHash.Hex(),
		"ETH",
		"mainnet",
		"transfer",
		log.Address,
		from,
		to,
		nftID.Uint64(),
	), nil
}

func (c *Contract) parseOrdersMatchedEvent(abi abi.ABI, log types.Log) (events.Event, error) {
	data, err := abi.Unpack(TopicOrdersMatched, log.Data)
	if err != nil {
		c.log.Error().Err(err).Str("topic", TopicOrdersMatched).Msg("Failed to unpack event")
		return nil, err
	}

	buyHash, ok := data[0].([32]uint8)
	if !ok {
		return nil, errors.New("invalid event format: missing buy hash")
	}

	sellHash, ok := data[1].([32]uint8)
	if !ok {
		return nil, errors.New("invalid event format: missing sell hash")
	}

	price, ok := data[2].(*big.Int)
	if !ok {
		return nil, errors.New("invalid event format: missing price")
	}

	return events.NewOrdersMatched(
		log.TxHash.Hex(),
		"ETH",
		"mainnet",
		"orders_matched",
		log.Address,
		common.BytesToHash(buyHash[:]),
		common.BytesToHash(sellHash[:]),
		common.BytesToAddress(log.Topics[1].Bytes()),
		common.BytesToAddress(log.Topics[2].Bytes()),
		price.Uint64(),
		common.BytesToHash(log.Topics[3].Bytes()),
	), nil
}
