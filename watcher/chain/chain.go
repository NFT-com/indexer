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

	return &w, nil
}

func (j *Watcher) Watch(_ context.Context) error {

	latestBlock := <-j.blocks
	j.latestBlock = latestBlock

	err := j.bootstrap()
	if err != nil {
		return fmt.Errorf("could not bootstrap system: %w", err)
	}

	for {
		select {

		case <-j.close:
			j.log.Debug().Msg("watching aborted")
			return nil

		case block := <-j.blocks:

			jobs := make([]jobs.Parsing, 0, len(j.config.Contracts))
			for _, contract := range j.config.Contracts {
				jobs = append(jobs, j.createJobsForContract(contract, block)...)
			}

			j.log.Debug().
				Uint64("height", block.Uint64()).
				Int("jobs", len(jobs)).
				Msg("processing block")

			err = j.publishJobs(jobs)
			if err != nil {
				j.log.Error().Err(err).Str("block", block.String()).Msg("could not create parsing jobs for block")
				continue
			}
		}
	}
}

func (j *Watcher) createJobsForContract(contract string, block *big.Int) []jobs.Parsing {
	contractIndexes, ok := j.config.StartIndexes[contract]
	if !ok {
		j.log.Error().Str("contract", contract).Msg("could not get contract indexes")
		return nil
	}

	standards, ok := j.config.Standards[contract]
	if !ok {
		j.log.Error().Str("contract", contract).Msg("could not get standards for contract")
		return nil
	}

	jobsList := make([]jobs.Parsing, 0, len(standards))
	for _, standard := range standards {
		collectionIndexes, ok := contractIndexes[standard]
		if !ok {
			j.log.Error().Str("standard", standard).Str("contract", contract).Msg("could not get standard indexes")
			continue
		}

		eTypes, ok := j.config.EventTypes[standard]
		if !ok {
			j.log.Error().Str("standard", standard).Str("contract", contract).Msg("could not get event types for standard")
			continue
		}

		for _, eType := range eTypes {
			startingBlock, ok := collectionIndexes[eType]
			if !ok {
				j.log.Error().Str("contract", contract).Str("standard", standard).Str("event_type", eType).Msg("could not check event type starting block")
				continue
			}

			// means that the current block index is lower that the starting block for this contract
			if block.CmpAbs(startingBlock) < 0 {
				continue
			}

			jobsList = append(jobsList, jobs.Parsing{
				ChainURL:     j.config.ChainURL,
				ChainID:      j.config.ChainID,
				ChainType:    j.config.ChainType,
				BlockNumber:  block.String(),
				Address:      contract,
				StandardType: standard,
				EventType:    eType,
			})
		}
	}

	return jobsList
}

func (j *Watcher) publishJobs(jobs []jobs.Parsing) error {
	batches := int64(len(jobs)) / j.config.Batch
	if int64(len(jobs))%j.config.Batch != 0 {
		batches++
	}

	for i := int64(0); i < batches; i++ {
		startBatch := i * j.config.Batch
		endBatch := (i + 1) * j.config.Batch
		if endBatch > int64(len(jobs)) {
			endBatch = int64(len(jobs))
		}

		batch := jobs[startBatch:endBatch]
		err := j.apiClient.CreateParsingJobs(batch)
		if err != nil {
			return fmt.Errorf("could not create parsing jobs: %w", err)
		}
	}

	return nil
}

func (j *Watcher) Close() {
	close(j.close)
}
