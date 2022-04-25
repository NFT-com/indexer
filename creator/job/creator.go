package job

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/creator"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/postgres"
)

type Creator struct {
	mutex    sync.Mutex
	log      zerolog.Logger
	start    uint64
	last     uint64
	template jobs.Parsing
	persist  creator.Persister
	check    creator.Checker
	limit    uint
}

func NewCreator(log zerolog.Logger, start uint64, template jobs.Parsing, persist creator.Persister, check creator.Checker, limit uint) *Creator {

	log = log.With().
		Str("component", "jobs_creator").
		Str("address", template.Address).
		Str("standard", template.Standard).
		Str("event", template.Event).
		Logger()

	c := Creator{
		mutex:    sync.Mutex{},
		log:      log,
		start:    start,
		last:     0,
		template: template,
		persist:  persist,
		check:    check,
		limit:    limit,
	}

	return &c
}

func (c *Creator) Notify(height uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if height < c.start {
		c.log.Debug().Uint64("height", height).Uint64("start", c.start).Msg("ignoring notify (start not reached)")
		return
	}

	if height <= c.last {
		c.log.Debug().Uint64("height", height).Uint64("last", c.last).Msg("ignoring notify (height stale)")
		return
	}

	// Check how many pending jobs are still in the DB for this combination and skip if there are too many.
	count, err := c.check.CountPendingParsingJobs(c.template.ChainURL, c.template.ChainType, c.template.Address, c.template.Standard, c.template.Event)
	if err != nil {
		c.log.Error().Err(err).Msg("could not count pending parsing jobs")
		return
	}

	if count >= c.limit {
		c.log.Debug().Uint("count", count).Uint("limit", c.limit).Msg("ignoring notify (limit reached")
		return
	}

	// If we have an existing job, we start at the height just above.
	start := uint64(0)
	existing, err := c.check.LastParsingJob(c.template.ChainID, c.template.Address, c.template.Event)
	if err != nil && !errors.Is(err, postgres.ErrResourceNotFound) {
		c.log.Error().Err(err).Msg("could not check on most recent job")
		return
	}
	if err == nil {
		start = existing.BlockNumber + 1
	}
	if start < c.start {
		start = c.start
	}
	if start > height {
		c.log.Debug().Uint64("height", height).Uint64("start", start).Msg("ignoring notify (all jobs pending)")
		return
	}

	created := uint(0)
	maximum := c.limit - count
	for index := start; index < start+uint64(maximum) && index <= height; index++ {
		job := c.template
		job.ID = uuid.New().String()
		job.BlockNumber = index
		c.persist.Store(&job)
		created++
		c.last = index
	}

	c.log.Info().Uint("jobs", created).Uint64("start", start).Uint64("end", c.last).Msg("created parsing jobs")
}
