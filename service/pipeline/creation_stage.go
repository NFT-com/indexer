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
	parsings     Publisher
	cfg          CreationConfig
}

func NewCreationStage(log zerolog.Logger, collections CollectionStore, marketplaces MarketplaceStore, boundaries BoundaryStore, parsings Publisher, options ...Option) *CreationStage {

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
		parsings:     parsings,
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
	var sentinel string
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
				Msg("updating start height with latest height")
		}
	}

	// We then enter a loop where we keep creating jobs until we hit the stop condition...
	created := uint(0)
	for created < c.cfg.BatchSize {

		// After determining the start height for every combination, we identify one of the
		// contract addresses with the lowest start height. We will limit the jobs to the
		// event hashes of that contract.
		lowest := uint64(math.MaxUint64)
		for _, combination := range combinations {
			if combination.StartHeight < lowest {
				lowest = combination.StartHeight
				sentinel = combination.ContractAddress
				c.log.Debug().
					Uint64("height", combination.StartHeight).
					Str("contract_address", combination.ContractAddress).
					Msg("updated sentinel smart contract for event hashes")
			}
		}

		// Next, we gather all event types for the given sentinel address. This step is
		// needed because we have split everything into combinations per event hash, so
		// we just match all of those with the same address here.
		hashSet := make(map[string]struct{})
		for _, combination := range combinations {
			if combination.ContractAddress == sentinel {
				hashSet[combination.EventHash] = struct{}{}
				c.log.Debug().
					Str("contract_address", combination.ContractAddress).
					Str("event_hash", combination.EventHash).
					Msg("added event hash for sentinel smart contract")
			}
		}

		// We start at the lowest start height.
		start := lowest

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

		// Now we want to include all of the contract addresses that have a start height
		// at or below our end height, and an event type that is part of the current run.
		addressSet := make(map[string]struct{})
		for _, combination := range combinations {

			// Stop if we have reached the maximum number of addresses.
			if uint(len(addressSet)) > c.cfg.AddressLimit {
				break
			}

			// Skip if start height above end.
			if combination.StartHeight > end {
				continue
			}

			// Skip if event hash is not in current hash set.
			_, ok := hashSet[combination.EventHash]
			if !ok {
				continue
			}

			addressSet[combination.ContractAddress] = struct{}{}
			combination.StartHeight = end + 1
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
		err = c.parsings.Publish(params.TopicParsing, payload)
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
