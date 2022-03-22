package web3

import (
	"context"
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/events"
)

const (
	indexBase = 10
)

type Web3 struct {
	ethClient *ethclient.Client
	chainID   string
	networkID string
}

func NewWeb3(ctx context.Context, url string) (*Web3, error) {
	ethClient, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	networkID, err := ethClient.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	w := Web3{
		ethClient: ethClient,
		chainID:   chainID.String(),
		networkID: networkID.String(),
	}

	return &w, nil
}

func (w *Web3) BlockEvents(ctx context.Context, blockNumber, eventType, contract string) ([]events.RawEvent, error) {
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
		return nil, err
	}

	evnts := make([]events.RawEvent, 0)
	for _, log := range logs {
		if log.Removed {
			continue
		}

		eventJson, err := log.MarshalJSON()
		if err != nil {
			continue
		}

		hash := sha256.Sum256(eventJson)

		indexData := make([]string, 0)
		for _, topic := range log.Topics[1:] {
			indexData = append(indexData, topic.String())
		}

		e := events.RawEvent{
			ID:              common.Bytes2Hex(hash[:]),
			ChainID:         w.chainID,
			NetworkID:       w.networkID,
			BlockNumber:     blockNumber,
			BlockHash:       log.BlockHash.String(),
			Address:         contract,
			TransactionHash: log.TxHash.String(),
			EventType:       log.Topics[0].String(),
			IndexData:       indexData,
			Data:            log.Data,
		}

		evnts = append(evnts, e)
	}

	return evnts, nil
}

func (w *Web3) Close() {
	w.ethClient.Close()
}
