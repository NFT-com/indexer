package database

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/persister"
)

type Persister struct {
	mutex     sync.Mutex
	log       zerolog.Logger
	ctx       context.Context
	tick      *time.Ticker
	store     persister.Store
	jobs      []*jobs.Parsing
	watermark uint
}

func NewPersister(log zerolog.Logger, ctx context.Context, store persister.Store, delay time.Duration, watermark uint) *Persister {

	p := Persister{
		mutex:     sync.Mutex{},
		log:       log,
		ctx:       ctx,
		tick:      time.NewTicker(delay),
		store:     store,
		jobs:      []*jobs.Parsing{},
		watermark: 1000,
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

	go p.execute()
}

func (p *Persister) execute() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.jobs) == 0 {
		return
	}

	err := p.store.CreateParsingJobs(p.jobs)
	if err != nil {
		p.log.Error().Err(err).Msg("could not create parsing jobs")
	}

	p.jobs = []*jobs.Parsing{}
}

func (p *Persister) process() {

ProcessLoop:
	for {

		select {

		case <-p.ctx.Done():

			p.log.Debug().Msg("terminating jobs persister")

			break ProcessLoop

		case <-p.tick.C:

			p.execute()
		}
	}

	p.tick.Stop()
}
