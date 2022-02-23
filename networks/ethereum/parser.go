package ethereum

import (
	"context"
	"crypto/sha256"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
)

type Parser struct {
	log zerolog.Logger

	network   string
	networkID string
	chainID   string
	client    Client
}

func NewParser(ctx context.Context, log zerolog.Logger, client Client) (*Parser, error) {
	if client == nil {
		return nil, errors.New("invalid ethereum client")
	}

	p := Parser{
		log:    log.With().Str("component", "parser").Logger(),
		client: client,
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	p.networkID = networkID.String()
	p.chainID = chainID.String()

	return &p, nil
}

func (p *Parser) Parse(ctx context.Context, block *block.Block) ([]*event.Event, error) {
	blockHash := common.HexToHash(block.String())

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
