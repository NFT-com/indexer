package persister

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/models/jobs"
)

type Persister struct {
	mutex     sync.Mutex
	log       zerolog.Logger
	ctx       context.Context
	tick      *time.Ticker
	parsings  ParsingStore
	jobs      []*jobs.Parsing
	watermark uint
}

func New(log zerolog.Logger, ctx context.Context, parsings ParsingStore, delay time.Duration, watermark uint) *Persister {

	p := Persister{
		mutex:     sync.Mutex{},
		log:       log,
		ctx:       ctx,
		tick:      time.NewTicker(delay),
		parsings:  parsings,
		jobs:      make([]*jobs.Parsing, 0, watermark),
		watermark: watermark,
	}

	go p.process()

	return &p
}

func (p *Persister) Store(job *jobs.Parsing) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.jobs = append(p.jobs, job)

	if uint(len(p.jobs)) < p.watermark {
		return
	}

	p.execute()
}

func (p *Persister) check() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.jobs) == 0 {
		return
	}

	p.execute()
}

func (p *Persister) execute() {

	err := p.parsings.Insert(p.jobs...)
	if err != nil {
		p.log.Error().Err(err).Msg("could not create parsing jobs")
	}

	p.log.Info().Int("jobs", len(p.jobs)).Msg("persisted parsing jobs")

	p.jobs = make([]*jobs.Parsing, 0, p.watermark)
}

func (p *Persister) process() {
	defer p.tick.Stop()

	for {
		select {

		case <-p.ctx.Done():
			p.log.Debug().Msg("terminating jobs persister")
			return

		case <-p.tick.C:

			p.check()
		}
	}
}
