package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/contracts"
	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/parse"
)

type Parser struct {
	log zerolog.Logger

	client   *ethclient.Client
	contract *contracts.Contract
}

func NewParser(log zerolog.Logger, client *ethclient.Client, contract *contracts.Contract) *Parser {
	p := Parser{
		log:      log.With().Str("component", "parser").Logger(),
		client:   client,
		contract: contract,
	}

	return &p
}

func (p *Parser) Parse(ctx context.Context, block *parse.Block) ([]events.Event, error) {
	hash := common.HexToHash(block.String())
	logger := p.log.With().Str("block_hash", block.String()).Logger()
	outEvents := make([]events.Event, 0, 64)

	query := ethereum.FilterQuery{
		BlockHash: &hash,
		Topics: [][]common.Hash{
			{
				contracts.TopicHash(contracts.TopicTransfer),
				contracts.TopicHash(contracts.TopicTransferSingle),
				contracts.TopicHash(contracts.TopicTransferBatch),
				contracts.TopicHash(contracts.TopicURI),
			},
		},
	}
	logs, err := p.client.FilterLogs(ctx, query)
	if err != nil {
		// FIXME: Should we stop the subscriber in this case?
		logger.Error().Err(err).Msg("could not filter ethereum client logs")
	}

	for _, l := range logs {
		if len(l.Topics) == 0 {
			logger.Warn().Msg("unexpected event: missing topic")
			continue
		}

		event, err := p.contract.ParseEvent(ctx, l)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("address", l.Address.Hex()).
				Str("topic", l.Topics[0].String()).
				Msg("could not parse event")
			continue
		}

		if event == nil {
			continue
		}

		outEvents = append(outEvents, event)
	}

	return outEvents, nil
}
