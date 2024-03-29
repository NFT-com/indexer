package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/NFT-com/indexer/config/retry"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlocksNotifier struct {
	log    zerolog.Logger
	ctx    context.Context
	wsURL  string
	heads  chan *types.Header
	sub    ethereum.Subscription
	done   chan struct{}
	errors chan error
	listen Listener
}

func NewBlocksNotifier(log zerolog.Logger, ctx context.Context, wsURL string, listen Listener) (*BlocksNotifier, error) {

	n := BlocksNotifier{
		log:    log.With().Str("component", "blocks_notifier").Logger(),
		ctx:    ctx,
		wsURL:  wsURL,
		listen: listen,
		heads:  make(chan *types.Header, 1),
		done:   make(chan struct{}),
		errors: make(chan error, 1),
	}

	n.subscribe(ctx)

	go n.process()

	return &n, nil
}

func (n *BlocksNotifier) subscribe(ctx context.Context) {

	var sub ethereum.Subscription

	notify := func(err error, dur time.Duration) {
		log.Warn().Err(err).Dur("duration", dur).Msg("could not subscribe to headers, retrying")
	}

	err := backoff.RetryNotify(func() error {

		select {
		case <-n.done:
			return nil
		default:
		}

		cli, err := ethclient.DialContext(ctx, n.wsURL)
		if err != nil {
			return fmt.Errorf("could not dial to websocket node api: %w", err)
		}

		sub, err = cli.SubscribeNewHead(ctx, n.heads)
		if err != nil {
			return fmt.Errorf("could not subscribe to heads: %w", err)
		}

		return nil
	}, retry.Indefinite(), notify)
	if err != nil {
		n.errors <- err
		return
	}

	n.sub = sub
}

func (n *BlocksNotifier) process() {

ProcessLoop:
	for {

		select {

		case <-n.ctx.Done():

			n.log.Debug().Msg("terminating blocks notifier")

			close(n.done)

			break ProcessLoop

		case err := <-n.errors:

			n.log.Error().Err(err).Msg("aborting blocks notifier")

			break ProcessLoop

		case err := <-n.sub.Err():

			n.log.Error().Err(err).Msg("error from websocket connection")

			go n.subscribe(n.ctx)

			continue ProcessLoop

		case head := <-n.heads:

			height := head.Number.Uint64()

			n.log.Debug().Uint64("height", height).Msg("notifying new block height")

			n.listen.Notify(height)

			continue ProcessLoop

		}

	}

	n.sub.Unsubscribe()
}
