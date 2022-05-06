package creator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/subchen/go-trylock/v2"
)

type Creator struct {
	mutex       trylock.TryLocker
	log         zerolog.Logger
	collections CollectionStore
	parsings    ParsingStore
	cfg         Config
}

func New(log zerolog.Logger, collections CollectionStore, parsings ParsingStore, options ...Option) *Creator {

	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	c := Creator{
		mutex:       trylock.New(),
		log:         log.With().Str("component", "jobs_creator").Logger(),
		collections: collections,
		parsings:    parsings,
		cfg:         cfg,
	}

	return &c
}

func (c *Creator) Notify(height uint64) {
	if !c.mutex.TryLock(context.Background()) {
		return
	}
	defer c.mutex.Unlock()

	err := c.execute(height)
	if err != nil {
		c.log.Error().Err(err).Msg("could not execute job creation")
		return
	}
}

func (c *Creator) execute(height uint64) error {

	// First, we get the number of pending jobs in the DB, so that we don't create
	// new jobs if we are above that.
	pending, err := c.parsings.Pending(c.cfg.ChainID)
	if err != nil {
		return fmt.Errorf("could not count pending jobs: %w", err)
	}

	// At this point, we can already bail if we have reached the limit of jobs.
	if pending >= c.cfg.PendingLimit {
		return nil
	}

	// Our goal now is to build a list of all possible combinations of contract
	// address and event hash for this chain.
	combinations, err := c.collections.Combinations(c.cfg.ChainID)
	if err != nil {
		return fmt.Errorf("could not get combinations: %w", err)
	}

	// Then, we get the latest job for each combination in order to update the
	// start height where necessary.
	for _, combination := range combinations {
		latest, err := c.parsings.Latest(combination.ChainID, combination.ContractAddress, combination.EventHash)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return fmt.Errorf("could not get latest parsing job: %w", err)
		}
		if latest >= combination.StartHeight {
			combination.StartHeight = latest + 1
		}
	}

	// We then enter a loop where we keep creating jobs until we hit the stop condition...
	for pending < c.cfg.PendingLimit {

		// First, go through all combinations and found the lowest start height, which
		// we will use as the start height for the next job we are creating.
		start := uint64(math.MaxUint64)
		for _, combination := range combinations {
			if combination.StartHeight < start {
				start = combination.StartHeight
			}
		}

		// The end height will be the lower between our configured height range and
		// the height that is available.
		end := start + uint64(c.cfg.HeightRange)
		if end > height {
			end = height
		}
		if end < start {
			c.log.Debug().Uint64("start", start).Uint64("end", end).Msg("no jobs to be created")
			break
		}

		// Now we want to include all of the contract addresses and all of the event hashes
		// that have a start height at or below our new start height.
		addressSet := make(map[string]struct{})
		hashSet := make(map[string]struct{})
		for _, combination := range combinations {
			if combination.StartHeight <= start {
				addressSet[combination.ContractAddress] = struct{}{}
				hashSet[combination.EventHash] = struct{}{}
				combination.StartHeight = end + 1
			}
		}

		// Now, we simply need to create the next job with the list of contract addresses
		// and event hashes and insert it into the database.
		addresses := make([]string, 0, len(addressSet))
		for address := range addressSet {
			addresses = append(addresses, address)
		}
		hashes := make([]string, 0, len(hashSet))
		for hash := range hashSet {
			hashes = append(hashes, hash)
		}
		parsing := jobs.Parsing{
			ID:                uuid.New().String(),
			ChainID:           c.cfg.ChainID,
			Status:            jobs.StatusCreated,
			ContractAddresses: addresses,
			EventHashes:       hashes,
			StartHeight:       start,
			EndHeight:         end,
		}
		err = c.parsings.Insert(&parsing)
		if err != nil {
			return fmt.Errorf("could not insert parsing job: %w", err)
		}

		pending++
	}

	return nil
}
