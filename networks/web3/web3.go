package web3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/log"
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
	return w.chainID, nil
}

func (w *Web3) SubscribeToBlocks(ctx context.Context, blocks chan *big.Int) error {
	headerChannel := make(chan *types.Header)
	subscription, err := w.ethClient.SubscribeNewHead(ctx, headerChannel)
	if err != nil {
		return fmt.Errorf("could not subscribe to new headers: %w", err)
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
		return nil, fmt.Errorf("could not get header: %w", err)
	}

	return header.Number, nil
}

func (w *Web3) BlockEvents(ctx context.Context, startBlock uint64, endBlock uint64, eventTypes []string, contracts []string) ([]log.RawLog, error) {
	start := big.NewInt(0).SetUint64(startBlock)
	end := big.NewInt(0).SetUint64(endBlock)

	addresses := make([]common.Address, 0, len(contracts))
	for _, contract := range contracts {
		addresses = append(addresses, common.HexToAddress(contract))
	}

	topics := make([][]common.Hash, 0, len(eventTypes))
	for _, event := range eventTypes {
		topics = append(topics, []common.Hash{common.HexToHash(event)})
	}

	query := ethereum.FilterQuery{
		FromBlock: start,
		ToBlock:   end,
		Addresses: addresses,
		Topics:    topics,
	}

	web3Logs, err := w.ethClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get filtered logs: %w", err)
	}

	logs := make([]log.RawLog, 0, len(web3Logs))
	for _, web3Log := range web3Logs {
		// in case that the transaction got reverted
		if web3Log.Removed {
			continue
		}

		header, err := w.ethClient.HeaderByNumber(ctx, big.NewInt(0).SetUint64(web3Log.BlockNumber))
		if err != nil {
			return nil, fmt.Errorf("could not get block header: %w", err)
		}

		blockDate := time.Unix(int64(header.Time), 0)

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
			BlockNumber:     header.Number.Uint64(),
			Index:           web3Log.Index,
			BlockHash:       web3Log.BlockHash.String(),
			Address:         web3Log.Address.String(),
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

func (w *Web3) CallContract(ctx context.Context, block *big.Int, sender, contract string, input []byte) ([]byte, error) {
	var (
		from    = common.HexToAddress(sender)
		address = common.HexToAddress(contract)
	)

	msg := ethereum.CallMsg{From: from, To: &address, Data: input}
	output, err := w.ethClient.CallContract(ctx, msg, block)
	if err != nil {
		return nil, fmt.Errorf("could not call contract: %w", err)
	}

	return output, nil
}

func (w *Web3) Close() {
	close(w.close)
	w.ethClient.Close()
}
