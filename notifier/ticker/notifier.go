package ticker

import (
	"context"
	"fmt"
	"time"

	"github.com/NFT-com/indexer/notifier"
	"github.com/rs/zerolog"
)

type Notifier struct {
	log     zerolog.Logger
	ctx     context.Context
	ticker  *time.Ticker
	heights chan uint64
	latest  uint64
	listen  notifier.Listener
}

func NewNotifier(log zerolog.Logger, ctx context.Context, period time.Duration, latest uint64, listen notifier.Listener) (*Notifier, error) {

	if period <= time.Millisecond {
		return nil, fmt.Errorf("invalid ticker period (%s)", period)
	}

	n := Notifier{
		log:     log,
		ctx:     ctx,
		ticker:  time.NewTicker(period),
		heights: make(chan uint64, 1),
		latest:  latest,
		listen:  listen,
	}

	go n.process()

	return &n, nil

}

func (n *Notifier) Notify(height uint64) {
	n.heights <- height
}

func (n *Notifier) process() {

ProcessLoop:
	for {

		select {

		case <-n.ctx.Done():

			n.log.Debug().Msg("terminating ticker notifications")

			break ProcessLoop

		case height := <-n.heights:

			n.log.Debug().Uint64("height", height).Msg("updating ticker height")

			n.latest = height

		case <-n.ticker.C:

			n.log.Debug().Uint64("height", n.latest).Msg("notifying ticker height")

			n.listen.Notify(n.latest)
		}
	}

	n.ticker.Stop()
	close(n.heights)
}
