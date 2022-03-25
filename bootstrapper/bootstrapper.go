package bootstrapper

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/client"
)

type Bootstrapper struct {
	log          zerolog.Logger
	apiClient    *client.Client
	chainURL     string
	chainType    string
	standardType string
	contract     string
	eventType    string
	startIndex   int64
	endIndex     int64
	close        chan struct{}
}

func New(
	log zerolog.Logger,
	apiClient *client.Client,
	chainURL, chainType, standardType, contract, eventType string,
	startIndex, endIndex int64,
) *Bootstrapper {
	b := Bootstrapper{
		log:          log.With().Str("component", "bootstrapper").Logger(),
		apiClient:    apiClient,
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

			job := jobs.Parsing{
				ChainURL:     b.chainURL,
				ChainType:    b.chainType,
				BlockNumber:  big.NewInt(index).String(),
				Address:      b.contract,
				StandardType: b.standardType,
				EventType:    b.eventType,
			}

			_, err := b.apiClient.CreateParsingJob(job)
			if err != nil {
				return fmt.Errorf("could not create parsing job for block %s: %w", job.BlockNumber, err)
			}

			index++
		}
	}
}

func (b *Bootstrapper) Close() {
	close(b.close)
}
