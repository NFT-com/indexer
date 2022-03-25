package chain

import (
	"context"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/service/client"
)

type Watcher struct {
	log          zerolog.Logger
	apiClient    *client.Client
	network      networks.Network
	chainURL     string
	chainType    string
	standardType string
	contract     string
	eventType    string
	startIndex   *big.Int
	blocks       chan *big.Int
	close        chan struct{}
}

func NewWatcher(
	log zerolog.Logger,
	ctx context.Context,
	apiClient *client.Client,
	network networks.Network,
	chainURL, chainType, standardType, contract, eventType string,
) (*Watcher, error) {
	w := Watcher{
		log:          log.With().Str("component", "watcher").Logger(),
		apiClient:    apiClient,
		network:      network,
		chainURL:     chainURL,
		chainType:    chainType,
		standardType: standardType,
		contract:     contract,
		eventType:    eventType,
		blocks:       make(chan *big.Int),
		close:        make(chan struct{}),
	}

	err := network.SubscribeToBlocks(ctx, w.blocks)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (j *Watcher) Watch(ctx context.Context) error {
	for {
		select {
		case block := <-j.blocks:
			job := jobs.Parsing{
				ChainURL:     j.chainURL,
				ChainType:    j.chainType,
				BlockNumber:  block.String(),
				Address:      j.contract,
				StandardType: j.standardType,
				EventType:    j.eventType,
			}

			_, err := j.apiClient.CreateParsingJob(job)
			if err != nil {
				j.log.Error().Err(err).Str("block", block.String()).Msg("failed create parsing job")
				continue
			}
		case <-j.close:
			return nil
		}
	}
}

func (j *Watcher) Close() {
	close(j.close)
}
