package bootstrapper

import (
	"context"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/service/client"
)

type Bootstrapper struct {
	log          zerolog.Logger
	apiClient    *client.Client
	network      networks.Network
	chainURL     string
	chainType    string
	standardType string
	contract     string
	eventType    string
	startIndex   int64
	endIndex     int64
	close        chan struct{}
}

func NewBootstrapper(
	log zerolog.Logger,
	apiClient *client.Client,
	network networks.Network,
	chainURL, chainType, standardType, contract, eventType string,
	startIndex, endIndex int64,
) *Bootstrapper {
	b := Bootstrapper{
		log:          log.With().Str("component", "bootstrapper").Logger(),
		apiClient:    apiClient,
		network:      network,
		chainURL:     chainURL,
		chainType:    chainType,
		standardType: standardType,
		contract:     contract,
		eventType:    eventType,
		startIndex:   startIndex,
		endIndex:     endIndex,
		close:        make(chan struct{}),
	}

	return &b
}

func (b *Bootstrapper) Bootstrap(ctx context.Context) error {
	index := b.startIndex
	for {
		select {
		case <-b.close:
			return nil
		default:
			if index > b.endIndex {
				return nil
			}

			err := b.publishParsingJob(ctx, big.NewInt(index))
			if err != nil {
				return err
			}

			index++
		}
	}
}

func (b *Bootstrapper) publishParsingJob(ctx context.Context, block *big.Int) error {
	events, err := b.network.BlockEvents(ctx, block.String(), b.eventType, b.contract)
	if err != nil {
		b.log.Error().Err(err).Str("block", block.String()).Msg("failed to get block events")
		return err
	}

	if len(events) == 0 {
		return nil
	}

	job := jobs.Parsing{
		ChainURL:     b.chainURL,
		ChainType:    b.chainType,
		BlockNumber:  block.String(),
		Address:      b.contract,
		StandardType: b.standardType,
		EventType:    b.eventType,
	}

	_, err = b.apiClient.CreateParsingJob(job)
	if err != nil {
		b.log.Error().Err(err).Str("block", block.String()).Msg("failed create parsing job")
		return err
	}

	return nil
}

func (b *Bootstrapper) Close() {
	close(b.close)
}
