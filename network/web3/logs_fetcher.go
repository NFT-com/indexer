package web3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LogsFetcher struct {
	client ethclient.Client
}

func NewLogsFetcher(client ethclient.Client) *LogsFetcher {

	l := LogsFetcher{
		client: client,
	}

	return &l
}

func (l *LogsFetcher) Fetch(addresses []string, eventTypes []string, from uint64, to uint64) ([]types.Log, error) {

	ethAddresses := make([]common.Address, 0, len(addresses))
	for _, address := range addresses {
		ethAddresses = append(ethAddresses, common.HexToAddress(address))
	}

	topics := make([]common.Hash, 0, len(eventTypes))
	for _, eventType := range eventTypes {
		topics = append(topics, common.HexToHash(eventType))
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(from),
		ToBlock:   big.NewInt(0).SetUint64(to),
		Addresses: ethAddresses,
		Topics:    [][]common.Hash{topics},
	}

	logs, err := l.client.FilterLogs(context.TODO(), query)
	if err != nil {
		return nil, fmt.Errorf("could not filter logs: %w", err)
	}

	return logs, nil
}
