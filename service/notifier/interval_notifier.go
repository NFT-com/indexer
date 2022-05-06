package notifier

import (
	"context"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
)

type IntervalNotifier struct {
	log     zerolog.Logger
	ctx     context.Context
	listen  Listener
	cfg     Config
	heights chan uint64
	latest  uint64
}

func NewIntervalNotifier(log zerolog.Logger, ctx context.Context, listen Listener, options ...Option) *IntervalNotifier {

	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	i := IntervalNotifier{
		log:     log,
		ctx:     ctx,
		listen:  listen,
		cfg:     cfg,
		heights: make(chan uint64, 1),
		latest:  0,
	}

	go i.process()

	return &i

}

func (i *IntervalNotifier) Notify(height uint64) {
	i.heights <- height
}

func (i *IntervalNotifier) process() {

	// Introduce a random jitter from [0-n.period) so that we don't hit the DB
	// for all combinations at the exact same time.
	delay := time.Duration(rand.Uint64() % uint64(i.cfg.NotifyInterval))
	time.Sleep(delay)

ProcessLoop:
	for {

		select {

		case <-i.ctx.Done():

			i.log.Debug().Msg("terminating ticker notifications")

			break ProcessLoop

		case height := <-i.heights:

			i.log.Debug().Uint64("height", height).Msg("updating ticker height")

			i.latest = height

			i.listen.Notify(height)

		case <-time.After(i.cfg.NotifyInterval):

			i.log.Debug().Uint64("height", i.latest).Msg("notifying ticker height")

			i.listen.Notify(i.latest)
		}
	}

	close(i.heights)
}
