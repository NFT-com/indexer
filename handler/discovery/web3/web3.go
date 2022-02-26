package web3

import (
	"context"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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
		eventJson, err := log.MarshalJSON()
		if err != nil {
			return err
		}

		hash := sha256.Sum256(eventJson)

		indexedData := log.Topics[1:]
		indexedDataString := make([]string, 0, len(indexedData))
		for _, data := range indexedData {
			indexedDataString = append(indexedDataString, data.String())
		}

		parseJob := queue.ParseJob{
			ID:              common.Bytes2Hex(hash[:]),
			NetworkID:       networkID.String(),
			ChainID:         chainID.String(),
			Block:           log.BlockNumber,
			TransactionHash: log.TxHash.String(),
			AddressType:     "erc721", // TODO: We code get from DB, if not present go get from network (code_at) and the publish in the DB
			Address:         log.Address.String(),
			Topic:           log.Topics[0].String(),
			IndexedData:     indexedDataString,
			Data:            log.Data,
		}

		err = w.prod.PublishParseJob(w.parseQueueName, parseJob)
		if err != nil {
			return err
		}
	}

	return nil
}
