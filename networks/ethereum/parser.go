package ethereum

import (
	"context"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
)

type Parser struct {
	log    zerolog.Logger
	client *ethclient.Client

	networkID       string
	chainID         string
	filterAddresses []common.Address
}

func NewParser(ctx context.Context, log zerolog.Logger, client *ethclient.Client, stringAddresses []string) (*Parser, error) {
	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	filterAddresses := make([]common.Address, 0, len(stringAddresses))
	for _, address := range stringAddresses {
		filterAddresses = append(filterAddresses, common.HexToAddress(address))
	}

	p := Parser{
		log:             log.With().Str("component", "parser").Logger(),
		client:          client,
		networkID:       networkID.String(),
		chainID:         chainID.String(),
		filterAddresses: filterAddresses,
	}

	return &p, nil
}

func (p *Parser) Parse(ctx context.Context, block *block.Block) ([]*event.Event, error) {
	blockHash := common.HexToHash(block.String())

	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: p.filterAddresses,
		Topics: [][]common.Hash{
			{
				TopicHash(TopicTransfer),
				TopicHash(TopicTransferSingle),
				TopicHash(TopicTransferBatch),
				TopicHash(TopicURI),
			},
		},
	}
	logs, err := p.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	evts := make([]*event.Event, 0, 64)
	for _, l := range logs {
		eventJson, err := l.MarshalJSON()
		if err != nil {
			p.log.Error().Err(err).Str("block_hash", block.String()).Msg("failed to marshal event")
			continue
		}

		hash := sha256.Sum256(eventJson)
		e := event.Event{
			ID:              common.Bytes2Hex(hash[:]),
			Network:         p.networkID,
			Chain:           p.chainID,
			Block:           l.BlockNumber,
			TransactionHash: l.TxHash,
			Address:         l.Address,
			Topic:           l.Topics[0],
			IndexedData:     l.Topics[1:],
			Data:            l.Data,
		}

		evts = append(evts, &e)
	}

	return evts, nil
}
