package notifier

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	maxReconnects = 9
)

type BlocksNotifier struct {
	log          zerolog.Logger
	ctx          context.Context
	websocketURL string
	heads        chan *types.Header
	sub          ethereum.Subscription
	listen       Listener
}

func NewBlocksNotifier(log zerolog.Logger, ctx context.Context, websocketURL string, listen Listener) (*BlocksNotifier, error) {

	n := BlocksNotifier{
		log:          log.With().Str("component", "blocks_notifier").Logger(),
		ctx:          ctx,
		websocketURL: websocketURL,
		heads:        make(chan *types.Header, 1),
		listen:       listen,
	}

	err := n.subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe to new headers: %w", err)
	}

	go n.process()

	return &n, nil
}

func (n *BlocksNotifier) subscribe(ctx context.Context) error {
	var cli *ethclient.Client

	err := backoff.Retry(func() error {
		var err error

		cli, err = ethclient.DialContext(ctx, n.websocketURL)
		if err != nil {
			return fmt.Errorf("could not dial to websocket node api: %w", err)
		}

		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxReconnects))
	if err != nil {
		return fmt.Errorf("could not connect to to node API: %w", err)
	}

	sub, err := cli.SubscribeNewHead(ctx, n.heads)
	if err != nil {
		return fmt.Errorf("could not subscribe to heads: %w", err)
	}

	n.sub = sub
	return nil
}

func (n *BlocksNotifier) process() {
ProcessLoop:
	for {

		select {

		case <-n.ctx.Done():

			n.log.Debug().Msg("terminating blocks notifier")

			break ProcessLoop

		case err := <-n.sub.Err():

			n.log.Error().Err(err).Msg("error from websocket connection")

			err = n.subscribe(n.ctx)
			if err != nil {
				n.log.Error().Err(err).Msg("error reconnection to websocket")
				break ProcessLoop
			}

			continue ProcessLoop

		case head := <-n.heads:

			height := head.Number.Uint64()

			n.log.Info().Uint64("height", height).Msg("notifying new block height")

			n.listen.Notify(height)

			continue ProcessLoop

		}

	}

	n.sub.Unsubscribe()
}
