package job

import (
	"errors"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/creator"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/postgres"
)

type Creator struct {
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
		Str("component", "jobs_creater").
		Str("address", template.Address).
		Str("standard", template.Standard).
		Str("event", template.Event).
		Logger()

	c := Creator{
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

	if height < c.start {
		c.log.Debug().Uint64("height", height).Uint64("start", c.start).Msg("ignoring notify (start not reached)")
		return
	}

	if height <= c.last {
		c.log.Debug().Uint64("height", height).Uint64("last", c.last).Msg("ignoring notify (height stale)")
		return
	}

	c.last = height

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
	existing, err := c.check.HighestBlockNumberParsingJob(c.template.ChainURL, c.template.ChainType, c.template.Address, c.template.Standard, c.template.Event)
	if err != nil && !errors.Is(err, postgres.ErrResourceNotFound) {
		c.log.Error().Err(err).Msg("could not check on most recent job")
		return
	}
	if err == nil {
		number, err := strconv.ParseUint(existing.BlockNumber, 10, 64)
		if err != nil {
			c.log.Error().Err(err).Str("number", existing.BlockNumber).Msg("could not parse block number")
			return
		}
		start = number + 1
	}
	if start > height {
		c.log.Debug().Uint64("height", height).Uint64("start", start).Msg("ignoring notify (all jobs pending)")
		return
	}

	created := 0
	maximum := c.limit - count
	for index := start; index <= start+uint64(maximum) && index <= height; index++ {
		job := c.template
		job.BlockNumber = strconv.FormatUint(index, 10)
		c.persist.Store(&job)
		created++
	}

	c.log.Info().Int("jobs", created).Msg("created parsing jobs")
}
