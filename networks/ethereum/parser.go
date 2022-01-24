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
	log zerolog.Logger

	network string
	chain   string
	client  *ethclient.Client
}

func NewParser(log zerolog.Logger, client *ethclient.Client, network string, chain string) *Parser {
	p := Parser{
		log:     log.With().Str("component", "parser").Logger(),
		network: network,
		chain:   chain,
		client:  client,
	}

	return &p
}

func (p *Parser) Parse(ctx context.Context, block *block.Block) ([]*events.Event, error) {
	blockHash := common.HexToHash(block.String())
	evts := make([]*events.Event, 0, 64)

	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
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

	for _, l := range logs {
		eventJson, err := l.MarshalJSON()
		if err != nil {
			p.log.Error().Err(err).Str("block_hash", block.String()).Msg("failed to marshal event")
			continue
		}

		hash := sha256.Sum256(eventJson)
		event := events.Event{
			ID:              common.Bytes2Hex(hash[:]),
			Chain:           p.network,
			Network:         p.chain,
			Block:           l.BlockNumber,
			TransactionHash: l.TxHash,
			Address:         l.Address,
			Topic:           l.Topics[0],
			IndexedData:     l.Topics[1:],
			Data:            l.Data,
		}

		evts = append(evts, &event)
	}

	return evts, nil
}
