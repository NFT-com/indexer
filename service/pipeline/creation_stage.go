package pipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/subchen/go-trylock/v2"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/models/jobs"
)

type CreationStage struct {
	mutex        trylock.TryLocker
	log          zerolog.Logger
	collections  CollectionStore
	marketplaces MarketplaceStore
	boundaries   BoundaryStore
	pub          Publisher
	cfg          CreationConfig
}

func NewCreationStage(log zerolog.Logger, collections CollectionStore, marketplaces MarketplaceStore, boundaries BoundaryStore, pub Publisher, options ...Option) *CreationStage {

	cfg := DefaultCreationConfig
	for _, option := range options {
		option(&cfg)
	}

	c := CreationStage{
		mutex:        trylock.New(),
		log:          log.With().Str("component", "jobs_creator").Logger(),
		collections:  collections,
		marketplaces: marketplaces,
		boundaries:   boundaries,
		pub:          pub,
		cfg:          cfg,
	}

	return &c
}

func (c *CreationStage) Notify(height uint64) {
	if !c.mutex.TryLock(context.Background()) {
		c.log.Debug().Msg("skipping job creation (already in progress)")
		return
	}
	defer c.mutex.Unlock()

	err := c.execute(height)
	if err != nil {
		c.log.Error().Err(err).Msg("could not execute job creation")
		return
	}
}

func (c *CreationStage) execute(height uint64) error {

	var combinations []*jobs.Combination

	// Build a list of all possible combinations of collection
	// and event hash for this chain.
	collectionCombinations, err := c.collections.Combinations(c.cfg.ChainID)
	if err != nil {
		return fmt.Errorf("could not get collection combinations: %w", err)
	}
	combinations = append(combinations, collectionCombinations...)

	// Build a list of all possible combinations of marketplace
	// address and event hash for this chain.
	marketplaceCombinations, err := c.marketplaces.Combinations(c.cfg.ChainID)
	if err != nil {
		return fmt.Errorf("could not get marketplace combinations: %w", err)
	}
	combinations = append(combinations, marketplaceCombinations...)

	// Then, we get the latest job for each combination in order to update the
	// start height where necessary.
	for _, combination := range combinations {
		last, err := c.boundaries.ForCombination(combination.ChainID, combination.ContractAddress, combination.EventHash)
		if errors.Is(err, sql.ErrNoRows) {
			c.log.Debug().
				Uint64("chain_id", combination.ChainID).
				Str("contract_address", combination.ContractAddress).
				Str("event_hash", combination.EventHash).
				Uint64("start_height", combination.StartHeight).
				Msg("no last job found, using start height")
			continue
		}
		if err != nil {
			return fmt.Errorf("could not get latest parsing job: %w", err)
		}
		if last >= combination.StartHeight {
			combination.StartHeight = last + 1
			c.log.Debug().
				Uint64("chain_id", combination.ChainID).
				Str("contract_address", combination.ContractAddress).
				Str("event_hash", combination.EventHash).
				Uint64("start_height", combination.StartHeight).
				Uint64("last_height", last).
				Msg("updating start height with latest heigth")
		}
	}

	// We then enter a loop where we keep creating jobs until we hit the stop condition...
	created := uint(0)
	for created < c.cfg.BatchSize {

		// First, go through all combinations and find the lowest start height, which
		// we will use as the start height for the next job we are creating.
		start := uint64(math.MaxUint64)
		for _, combination := range combinations {
			if combination.StartHeight < start {
				start = combination.StartHeight
			}
		}

		// The end height will be the lower between our configured height range and
		// the height that is available.
		end := start + uint64(c.cfg.HeightLimit) - 1
		if end > height {
			end = height
		}
		if end < start {
			c.log.Debug().Uint64("start", start).Uint64("end", end).Msg("skipping job creation (no jobs left)")
			break
		}

		// Now we want to include all of the contract addresses and all of the event hashes
		// that have a start height at or below our end height.
		addressSet := make(map[string]struct{})
		hashSet := make(map[string]struct{})
		for _, combination := range combinations {
			if combination.StartHeight <= end {
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
			ID:                uuid.NewString(),
			ChainID:           c.cfg.ChainID,
			ContractAddresses: addresses,
			EventHashes:       hashes,
			StartHeight:       start,
			EndHeight:         end,
		}
		payload, err := json.Marshal(parsing)
		if err != nil {
			return fmt.Errorf("could not encode parsing job: %w", err)
		}
		err = c.pub.Publish(params.TopicParsing, payload)
		if err != nil {
			return fmt.Errorf("could not insert parsing job: %w", err)
		}

		created++

		err = c.boundaries.Upsert(c.cfg.ChainID, parsing.ContractAddresses, parsing.EventHashes, parsing.EndHeight, parsing.ID)
		if err != nil {
			return fmt.Errorf("could not update combination boundaries: %w", err)
		}

		c.log.Info().
			Uint64("chain_id", parsing.ChainID).
			Strs("contract_addresses", parsing.ContractAddresses).
			Strs("event_hashes", parsing.EventHashes).
			Uint64("start_height", parsing.StartHeight).
			Uint64("end_height", parsing.EndHeight).
			Msg("parsing job published")
	}

	return nil
}
