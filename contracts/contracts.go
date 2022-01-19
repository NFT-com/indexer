package contracts

import (
	"bytes"
	"context"
	"encoding/binary"
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
	TopicTransfer       = "Transfer"
	TopicTransferSingle = "TransferSingle"
	TopicTransferBatch  = "TransferBatch"
	TopicURI            = "URI"
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

	data, err := parsedAbi.Unpack(event.Name, log.Data)
	if err != nil {
		return nil, err
	}

	switch event.Name {
	case TopicTransfer:
		return c.parseTransferEvent(data, log)
	case TopicTransferSingle:
		return c.parseTransferSingleEvent(data, log)
	case TopicTransferBatch:
		return c.parseTransferBatchEvent(data, log)
	case TopicURI:
		return c.parseURIEvent(data, log)
	default:
		// FIXME: Handle more topics.
	}

	return nil, nil
}

func (c *Contract) parseTransferEvent(data []interface{}, log types.Log) (events.Event, error) {
	from, ok := data[0].(common.Address)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	to, ok := data[1].(common.Address)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	id, ok := data[2].(*big.Int)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	e := events.Transfer{
		ID:      log.TxHash.Hex(),
		Chain:   "ETH",
		Network: "mainnet",
		Topic:   "transfer",
		Address: log.Address,
		From:    from,
		To:      to,
		NftID:   id.Uint64(),
	}

	return &e, nil
}

func (c *Contract) parseTransferSingleEvent(data []interface{}, log types.Log) (events.Event, error) {
	var (
		operator = common.BytesToAddress(log.Topics[1].Bytes())
		from     = common.BytesToAddress(log.Topics[2].Bytes())
		to       = common.BytesToAddress(log.Topics[3].Bytes())
	)

	// FIXME: This should be uint256 -> big ints not uint64
	id, ok := data[0].(*big.Int)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	value, ok := data[1].(*big.Int)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	e := events.TransferSingle{
		ID:       log.TxHash.Hex(),
		Chain:    "ETH",
		Network:  "mainnet",
		Topic:    "transfer",
		Address:  log.Address,
		Operator: operator,
		From:     from,
		To:       to,
		NftID:    id.Uint64(),
		Value:    value.Uint64(),
	}

	return &e, nil
}

func (c *Contract) parseTransferBatchEvent(data []interface{}, log types.Log) (events.Event, error) {
	var (
		operator = common.BytesToAddress(log.Topics[1].Bytes())
		from     = common.BytesToAddress(log.Topics[2].Bytes())
		to       = common.BytesToAddress(log.Topics[3].Bytes())
		ids      = make([]uint64, 0)
		values   = make([]uint64, 0)
	)

	bigIntIDs, ok := data[0].([]*big.Int)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	for _, v := range bigIntIDs {
		ids = append(ids, v.Uint64())
	}

	bigIntValues, ok := data[1].([]*big.Int)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	for _, v := range bigIntValues {
		values = append(values, v.Uint64())
	}

	e := events.TransferBatch{
		ID:       log.TxHash.Hex(),
		Chain:    "ETH",
		Network:  "mainnet",
		Topic:    "transfer",
		Address:  log.Address,
		Operator: operator,
		From:     from,
		To:       to,
		NftIDs:   ids,
		Values:   values,
	}

	return &e, nil
}

func (c *Contract) parseURIEvent(data []interface{}, log types.Log) (events.Event, error) {
	value, ok := data[0].(string)
	if !ok {
		return nil, errors.New("invalid event format")
	}

	// FIXME: Is Ethereum BigEndian or LittleEndian?
	// FIXME: Is this the better way to parse?
	var id uint64
	err := binary.Read(bytes.NewBuffer(log.Topics[1].Bytes()), binary.BigEndian, &id)
	if err != nil {
		return nil, errors.New("invalid event format")
	}

	e := events.URI{
		ID:      log.TxHash.Hex(),
		Chain:   "ETH",
		Network: "mainnet",
		Topic:   "transfer",
		Address: log.Address,
		NftID:   id,
		URI:     value,
	}

	return &e, nil
}
