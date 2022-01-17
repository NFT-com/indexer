package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/contracts"
	"github.com/NFT-com/indexer/events"
)

type Parser struct {
	log zerolog.Logger

	client   *ethclient.Client
	contract *contracts.Contract
}

func NewParser(log zerolog.Logger, client *ethclient.Client, contract *contracts.Contract) *Parser {
	return &Parser{
		log:      log.With().Str("component", "parser").Logger(),
		client:   client,
		contract: contract,
	}
}

func (p *Parser) ParseBlock(ctx context.Context, block *block.Block) ([]events.Event, error) {
	hash := common.HexToHash(block.Hash)
	outEvents := make([]events.Event, 0, 64)

	// FIXME: Filter by address/topic as well, to fetch only relevant data and speed up the parsing process.
	logs, err := p.client.FilterLogs(ctx, ethereum.FilterQuery{BlockHash: &hash})
	if err != nil {
		// FIXME: Should we stop the subscriber in this case?
		p.log.Error().Err(err).Msg("could not filter ethereum client logs")
	}

	for _, l := range logs {
		if len(l.Topics) == 0 {
			p.log.Error().Msg("unexpected event: missing topic")
			continue
		}

		event, err := p.contract.ParseEvent(ctx, l)
		if err != nil {
			p.log.Error().Err(err).Str("address", l.Address.Hex()).Msg("could not parse event")
			continue
		}

		if event == nil {
			continue
		}

		outEvents = append(outEvents, event)
	}

	return outEvents, nil
}
