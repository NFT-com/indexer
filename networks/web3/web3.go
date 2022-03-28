package web3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/log"
)

const (
	// represents the decimal base (10) to parse the string numbers into *big.Int
	indexBase = 10
)

type Web3 struct {
	ethClient *ethclient.Client
	chainID   string
	close     chan struct{}
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
		close:     make(chan struct{}),
	}

	return &w, nil
}

func (w *Web3) ChainID(ctx context.Context) (string, error) {
	chainID, err := w.ethClient.ChainID(ctx)
	if err != nil {
		return "", err
	}

	return chainID.String(), nil
}

func (w *Web3) SubscribeToBlocks(ctx context.Context, blocks chan *big.Int) error {
	headerChannel := make(chan *types.Header)
	subscription, err := w.ethClient.SubscribeNewHead(ctx, headerChannel)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case header := <-headerChannel:
				blocks <- header.Number
			case <-w.close:
				subscription.Unsubscribe()
			}
		}
	}()

	return nil
}

func (w *Web3) GetLatestBlockHeight(ctx context.Context) (*big.Int, error) {
	header, err := w.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	return header.Number, nil
}

func (w *Web3) BlockEvents(ctx context.Context, blockNumber, eventType, contract string) ([]log.RawLog, error) {
	index, ok := big.NewInt(0).SetString(blockNumber, indexBase)
	if !ok {
		return nil, fmt.Errorf("could not parse block number into big.Int")
	}

	query := ethereum.FilterQuery{
		FromBlock: index,
		ToBlock:   index,
		Addresses: []common.Address{common.HexToAddress(contract)},
		Topics:    [][]common.Hash{{common.HexToHash(eventType)}},
	}

	web3Logs, err := w.ethClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get filtered logs: %w", err)
	}

	header, err := w.ethClient.HeaderByNumber(ctx, index)
	if err != nil {
		return nil, fmt.Errorf("could not get block header: %w", err)
	}

	blockDate := time.Unix(int64(header.Time), 0)

	logs := make([]log.RawLog, 0, len(web3Logs))
	for _, web3Log := range web3Logs {
		if web3Log.Removed {
			continue
		}

		eventJson, err := web3Log.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("could not marshal events to json: %w", err)
		}

		hash := sha256.Sum256(eventJson)

		indexData := make([]string, 0, len(web3Log.Topics)-1)
		for _, topic := range web3Log.Topics[1:] {
			indexData = append(indexData, topic.String())
		}

		l := log.RawLog{
			ID:              common.Bytes2Hex(hash[:]),
			ChainID:         w.chainID,
			BlockNumber:     blockNumber,
			BlockHash:       web3Log.BlockHash.String(),
			Address:         contract,
			TransactionHash: web3Log.TxHash.String(),
			EventType:       web3Log.Topics[0].String(),
			IndexData:       indexData,
			Data:            web3Log.Data,
			EmittedAt:       blockDate,
		}

		logs = append(logs, l)
	}

	return logs, nil
}

func (w *Web3) Close() {
	close(w.close)
	w.ethClient.Close()
}
