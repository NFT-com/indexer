package ticker

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/NFT-com/indexer/notifier"
	"github.com/rs/zerolog"
)

type Notifier struct {
	log     zerolog.Logger
	ctx     context.Context
	period  time.Duration
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
		period:  period,
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

	// Introduce a random jitter from [0-n.period) so that we don't hit the DB
	// for all combinations at the exact same time.
	delay := time.Duration(rand.Uint64() % uint64(n.period))
	time.Sleep(delay)

	// Initialize the ticker with the given period after jitter delay.
	ticker := time.NewTicker(n.period)

ProcessLoop:
	for {

		select {

		case <-n.ctx.Done():

			n.log.Debug().Msg("terminating ticker notifications")

			break ProcessLoop

		case height := <-n.heights:

			n.log.Debug().Uint64("height", height).Msg("updating ticker height")

			n.latest = height

		case <-ticker.C:

			n.log.Debug().Uint64("height", n.latest).Msg("notifying ticker height")

			n.listen.Notify(n.latest)
		}
	}

	ticker.Stop()
	close(n.heights)
}
