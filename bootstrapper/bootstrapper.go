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
	log       zerolog.Logger
	apiClient *client.Client
	config    Config
	close     chan struct{}
}

func New(
	log zerolog.Logger,
	apiClient *client.Client,
	config Config,
) *Bootstrapper {
	b := Bootstrapper{
		log:       log.With().Str("component", "bootstrapper").Logger(),
		apiClient: apiClient,
		config:    config,
		close:     make(chan struct{}),
	}

	return &b
}

func (b *Bootstrapper) Bootstrap(ctx context.Context) error {
	index := b.config.StartIndex

	for {
		select {
		case <-b.close:
			return nil
		default:
		}

		if index > b.config.EndIndex {
			return nil
		}

		job := jobs.Parsing{
			ChainURL:     b.config.ChainURL,
			ChainType:    b.config.ChainType,
			BlockNumber:  big.NewInt(index).String(),
			Address:      b.config.Contract,
			StandardType: b.config.StandardType,
			EventType:    b.config.EventType,
		}

		_, err := b.apiClient.CreateParsingJob(job)
		if err != nil {
			return fmt.Errorf("could not create parsing job for block %s: %w", job.BlockNumber, err)
		}

		index++
	}
}

func (b *Bootstrapper) Close() {
	close(b.close)
}
