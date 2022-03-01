package web3

import (
	"context"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"math/big"

	"github.com/NFT-com/indexer/queue"
	"github.com/NFT-com/indexer/queue/producer"
)

const (
	IndexBase = 10
)

type Web3 struct {
	log            zerolog.Logger
	parseQueueName string
	prod           *producer.Producer
}

func NewWeb3(log zerolog.Logger, parseQueueName string, prod *producer.Producer) (*Web3, error) {
	w := Web3{
		log:            log.With().Str("component", "discovery_web3").Logger(),
		parseQueueName: parseQueueName,
		prod:           prod,
	}

	return &w, nil
}

func (w *Web3) Handle(ctx context.Context, job queue.DiscoveryJob) error {
	client, err := ethclient.Dial(job.ChainURL)
	if err != nil {
		return err
	}

	query := w.getFilterQuery(job)
	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return err
	}

	for _, log := range logs {
		contractType := "erc721"

		parseJob, err := w.parseLog(log, networkID.String(), chainID.String(), contractType)
		if err != nil {
			return err
		}

		err = w.prod.PublishParseJob(w.parseQueueName, parseJob)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Web3) getFilterQuery(job queue.DiscoveryJob) ethereum.FilterQuery {
	zero := big.NewInt(0)
	startIndex, _ := zero.SetString(job.StartIndex, IndexBase)
	endIndex, _ := zero.SetString(job.EndIndex, IndexBase)

	addresses := make([]common.Address, 0, len(job.Contracts))
	for _, contract := range job.Contracts {
		addresses = append(addresses, common.HexToAddress(contract))
	}

	query := ethereum.FilterQuery{
		FromBlock: startIndex,
		ToBlock:   endIndex,
		Addresses: addresses,
	}

	return query
}

func (w *Web3) parseLog(log types.Log, networkID string, chainID string, contractType string) (queue.ParseJob, error) {
	eventJson, err := log.MarshalJSON()
	if err != nil {
		return queue.ParseJob{}, err
	}

	hash := sha256.Sum256(eventJson)

	indexedData := log.Topics[1:]
	indexedDataString := make([]string, 0, len(indexedData))
	for _, data := range indexedData {
		indexedDataString = append(indexedDataString, data.String())
	}

	parseJob := queue.ParseJob{
		ID:              common.Bytes2Hex(hash[:]),
		NetworkID:       networkID,
		ChainID:         chainID,
		Block:           log.BlockNumber,
		TransactionHash: log.TxHash.String(),
		AddressType:     contractType,
		Address:         log.Address.String(),
		Topic:           log.Topics[0].String(),
		IndexedData:     indexedDataString,
		Data:            log.Data,
	}

	return parseJob, nil
}
