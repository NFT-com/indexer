package heads

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/service/notifier"
)

type Notifier struct {
	log    zerolog.Logger
	ctx    context.Context
	heads  chan *types.Header
	sub    ethereum.Subscription
	listen notifier.Listener
}

func NewNotifier(log zerolog.Logger, ctx context.Context, cli *ethclient.Client, listen notifier.Listener) (*Notifier, error) {

	heads := make(chan *types.Header, 1)
	sub, err := cli.SubscribeNewHead(ctx, heads)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe to heads: %w", err)
	}

	n := Notifier{
		log:    log.With().Str("component", "heads_notifier").Logger(),
		ctx:    ctx,
		heads:  heads,
		sub:    sub,
		listen: listen,
	}

	go n.process()

	return &n, nil
}

func (n *Notifier) process() {

ProcessLoop:
	for {

		select {

		case <-n.ctx.Done():

			n.log.Debug().Msg("terminating heads notifier")

			break ProcessLoop

		case err := <-n.sub.Err():

			n.log.Error().Err(err).Msg("aborting heads watcher")

			break ProcessLoop

		case head := <-n.heads:

			height := head.Number.Uint64()

			n.log.Info().Uint64("height", height).Msg("notifying new head")

			n.listen.Notify(height)

			continue ProcessLoop

		}

	}

	n.sub.Unsubscribe()
}
