package web3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/event"
)

const (
	// represents the decimal base (10) to parse the string numbers into *big.Int
	indexBase = 10
)

type Web3 struct {
	ethClient *ethclient.Client
	chainID   string
	networkID string
}

func New(ctx context.Context, url string) (*Web3, error) {
	ethClient, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("could not dial to web3 client: %w", err)
	}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get chain id: %w", err)
	}

	w := Web3{
		ethClient: ethClient,
		chainID:   chainID.String(),
	}

	return &w, nil
}

func (w *Web3) BlockEvents(ctx context.Context, blockNumber, eventType, contract string) ([]event.RawEvent, error) {
	zero := big.NewInt(0)
	startIndex, _ := zero.SetString(blockNumber, indexBase)
	endIndex, _ := zero.SetString(blockNumber, indexBase)

	query := ethereum.FilterQuery{
		FromBlock: startIndex,
		ToBlock:   endIndex,
		Addresses: []common.Address{common.HexToAddress(contract)},
		Topics:    [][]common.Hash{{common.HexToHash(eventType)}},
	}

	logs, err := w.ethClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get filtered logs: %w", err)
	}

	evts := make([]event.RawEvent, 0, len(logs))
	for _, log := range logs {
		if log.Removed {
			continue
		}

		eventJson, err := log.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("could not marshal event to json: %w", err)
		}

		hash := sha256.Sum256(eventJson)

		indexData := make([]string, 0, len(log.Topics)-1)
		for _, topic := range log.Topics[1:] {
			indexData = append(indexData, topic.String())
		}

		e := event.RawEvent{
			ID:              common.Bytes2Hex(hash[:]),
			ChainID:         w.chainID,
			BlockNumber:     blockNumber,
			BlockHash:       log.BlockHash.String(),
			Address:         contract,
			TransactionHash: log.TxHash.String(),
			EventType:       log.Topics[0].String(),
			IndexData:       indexData,
			Data:            log.Data,
		}

		evts = append(evts, e)
	}

	return evts, nil
}

func (w *Web3) Close() {
	w.ethClient.Close()
}
