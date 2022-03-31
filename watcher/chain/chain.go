package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/service/client"
)

type Watcher struct {
	log         zerolog.Logger
	apiClient   *client.Client
	network     networks.Network
	config      Config
	latestBlock *big.Int
	blocks      chan *big.Int
	close       chan struct{}
}

func NewWatcher(
	log zerolog.Logger,
	ctx context.Context,
	apiClient *client.Client,
	network networks.Network,
	config Config,
) (*Watcher, error) {
	w := Watcher{
		log:       log.With().Str("component", "watcher").Logger(),
		apiClient: apiClient,
		network:   network,
		config:    config,
		blocks:    make(chan *big.Int),
		close:     make(chan struct{}),
	}

	err := network.SubscribeToBlocks(ctx, w.blocks)
	if err != nil {
		return nil, err
	}

	latestBlock := <-w.blocks
	w.latestBlock = latestBlock

	return &w, nil
}

func (j *Watcher) Watch(_ context.Context) error {
	err := j.bootstrap()
	if err != nil {
		return fmt.Errorf("could not bootstrap system: %w", err)
	}

	for {
		select {
		case <-j.close:
			return nil
		case block := <-j.blocks:
			job := jobs.Parsing{
				ChainURL:     j.config.ChainURL,
				ChainType:    j.config.ChainType,
				BlockNumber:  block.String(),
				Address:      j.config.Contract,
				StandardType: j.config.StandardType,
				EventType:    j.config.EventType,
			}

			_, err := j.apiClient.CreateParsingJob(job)
			if err != nil {
				j.log.Error().Err(err).Str("block", block.String()).Msg("could not create parsing job for block")
				continue
			}
		}
	}
}

func (j *Watcher) Close() {
	close(j.close)
}

func (j *Watcher) bootstrap() error {
	startingBlock, ok := big.NewInt(0).SetString(j.config.StartIndex, 0)
	if !ok {
		return fmt.Errorf("could not parse block number into big.Int")
	}

	index := startingBlock
	for {
		select {
		case <-j.close:
			return nil
		default:
		}

		job := jobs.Parsing{
			ChainURL:     j.config.ChainURL,
			ChainType:    j.config.ChainType,
			BlockNumber:  index.String(),
			Address:      j.config.Contract,
			StandardType: j.config.StandardType,
			EventType:    j.config.EventType,
		}

		_, err := j.apiClient.CreateParsingJob(job)
		if err != nil {
			return fmt.Errorf("could not create parsing job for block %s: %w", job.BlockNumber, err)
		}

		index = index.Add(index, big.NewInt(1))
	}
}
