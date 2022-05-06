package notifier

import (
	"context"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
)

type TickerNotifier struct {
	log     zerolog.Logger
	ctx     context.Context
	listen  Listener
	cfg     Config
	heights chan uint64
	latest  uint64
}

func NewTickerNotifier(log zerolog.Logger, ctx context.Context, listen Listener, options ...Option) *TickerNotifier {

	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	t := TickerNotifier{
		log:     log.With().Str("component", "ticker_notifier").Logger(),
		ctx:     ctx,
		listen:  listen,
		cfg:     cfg,
		heights: make(chan uint64, 1),
		latest:  0,
	}

	go t.process()

	return &t

}

func (i *TickerNotifier) Notify(height uint64) {
	i.heights <- height
}

func (t *TickerNotifier) process() {

	// Introduce a random jitter from [0-n.period) so that we don't hit the DB
	// for all combinations at the exact same time.
	delay := time.Duration(rand.Uint64() % uint64(t.cfg.NotifyInterval))
	time.Sleep(delay)

ProcessLoop:
	for {

		select {

		case <-t.ctx.Done():

			t.log.Debug().Msg("terminating ticker notifications")

			break ProcessLoop

		case height := <-t.heights:

			t.log.Debug().Uint64("height", height).Msg("updating ticker height")

			t.latest = height

			t.listen.Notify(height)

		case <-time.After(t.cfg.NotifyInterval):

			t.log.Debug().Uint64("height", t.latest).Msg("notifying ticker height")

			t.listen.Notify(t.latest)
		}
	}

	close(t.heights)
}
