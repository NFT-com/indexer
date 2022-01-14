package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/contracts"
	"github.com/NFT-com/indexer/events"
)

type LiveSubscriber struct {
	log zerolog.Logger

	client   *ethclient.Client
	contract *contracts.Contract
}

func NewLive(log zerolog.Logger, client *ethclient.Client, contract *contracts.Contract) *LiveSubscriber {
	l := LiveSubscriber{
		log:      log.With().Str("component", "live_subscriber").Logger(),
		client:   client,
		contract: contract,
	}

	return &l
}

func (s *LiveSubscriber) Subscribe(ctx context.Context, events chan events.Event) error {
	headerChannel := make(chan *types.Header)
	sub, err := s.client.SubscribeNewHead(ctx, headerChannel)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case header := <-headerChannel:
				blockHash := header.Hash()

				// FIXME: Filter by address/topic as well, to fetch only relevant data and speed up the parsing process.
				// FIXME: Make this code DRY since it is the exact same code for both subscribers.
				logs, err := s.client.FilterLogs(ctx, ethereum.FilterQuery{BlockHash: &blockHash})
				if err != nil {
					// FIXME: Should we stop the subscriber in this case?
					s.log.Error().Err(err).Msg("could not filter ethereum client logs")
				}

				for _, l := range logs {
					if len(l.Topics) == 0 {
						s.log.Error().Msg("unexpected event: missing topic")
						continue
					}

					event, err := s.contract.ParseEvent(ctx, l)
					if err != nil {
						s.log.Error().Err(err).Str("address", l.Address.Hex()).Msg("could not parse event")
						continue
					}

					if event == nil {
						continue
					}

					events <- event
				}
			}
		}
	}()

	err = <-sub.Err()
	return err
}

// Status returns nil if the subscriber is successfully connected to the Ethereum network.
// Otherwise, it returns an error.
func (s *LiveSubscriber) Status(ctx context.Context) error {
	_, err := s.client.NetworkID(ctx)
	return err
}

func (s *LiveSubscriber) Close() error {
	s.client.Close()
	return nil
}

type HistoricalSubscriber struct {
	log zerolog.Logger

	client     *ethclient.Client
	contract   *contracts.Contract
	startIndex int64
	endIndex   int64
}

func NewHistorical(log zerolog.Logger, client *ethclient.Client, contract *contracts.Contract, startIndex int64, endIndex int64) *HistoricalSubscriber {
	h := HistoricalSubscriber{
		log:        log.With().Str("component", "historical_subscriber").Logger(),
		client:     client,
		contract:   contract,
		startIndex: startIndex,
		endIndex:   endIndex,
	}

	return &h
}

func (s *HistoricalSubscriber) Subscribe(ctx context.Context, events chan events.Event) error {
	for i := s.startIndex; i <= s.endIndex; i++ {
		header, err := s.client.HeaderByNumber(ctx, big.NewInt(i))
		if err != nil {
			return fmt.Errorf("height: %v error: %v", i, err)
		}

		blockHash := header.Hash()

		// FIXME: Filter by address/topic as well, to fetch only relevant data and speed up the parsing process.
		// FIXME: Make this code DRY since it is the exact same code for both subscribers.
		logs, err := s.client.FilterLogs(ctx, ethereum.FilterQuery{BlockHash: &blockHash})
		if err != nil {
			// FIXME: Should we stop the subscriber in this case?
			s.log.Error().Err(err).Msg("could not filter ethereum client logs")
		}

		for _, l := range logs {
			if len(l.Topics) == 0 {
				s.log.Error().Msg("unexpected event: missing topic")
				continue
			}

			event, err := s.contract.ParseEvent(ctx, l)
			if err != nil {
				s.log.Error().Err(err).Str("address", l.Address.Hex()).Msg("could not parse event")
				continue
			}

			if event == nil {
				continue
			}

			events <- event
		}
	}

	return nil
}

func (s *HistoricalSubscriber) Status(ctx context.Context) error {
	_, err := s.client.NetworkID(ctx)
	return err
}

func (s *HistoricalSubscriber) Close() error {
	s.client.Close()
	return nil
}
