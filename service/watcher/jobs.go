package watcher

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/service/pipeline"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

type Job struct {
	log      zerolog.Logger
	parsings ParsingStore
	actions  ActionStore
	creator  *pipeline.JobCreator
	delay    time.Duration
	close    chan struct{}
}

func New(log zerolog.Logger, parsings ParsingStore, actions ActionStore, creator *pipeline.JobCreator, delay time.Duration) *Job {

	j := Job{
		log:      log.With().Str("component", "jobs_watcher").Logger(),
		parsings: parsings,
		actions:  actions,
		creator:  creator,
		delay:    delay,
		close:    make(chan struct{}),
	}

	return &j
}

func (j *Job) Watch() {
	go j.watchParsings()
	go j.watchActions()
}

func (j *Job) Close() {
	close(j.close)
}

func (j *Job) watchParsings() {

	for {
		select {

		case <-time.After(j.delay):

			parsings, err := j.parsings.List(jobs.StatusCreated)
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve parsing jobs")
				continue
			}

			if len(parsings) == 0 {
				j.log.Trace().Msg("skipping parsing jobs queuing (no parsing jobs left)")
				continue
			}

			err = j.handleParsingJobs(parsings)
			if err != nil {
				j.log.Error().Err(err).Msg("could not handle parsing jobs")
				continue
			}

		case <-j.close:
			return
		}
	}
}

func (j *Job) watchActions() {

	for {
		select {

		case <-time.After(j.delay):

			actions, err := j.actions.List(jobs.StatusCreated)
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve action jobs")
				continue
			}
			if len(actions) == 0 {
				j.log.Trace().Msg("skipping action jobs queuing (no action jobs left)")
				continue
			}

			err = j.handleActionJobs(actions)
			if err != nil {
				j.log.Error().Err(err).Msg("could not handle action jobs")
				continue
			}

		case <-j.close:
			return
		}
	}
}

func (j *Job) handleParsingJobs(parsings []*jobs.Parsing) error {

	parsingIDs := make([]string, 0, len(parsings))
	for _, parsing := range parsings {
		parsingIDs = append(parsingIDs, parsing.ID)
	}

	err := j.parsings.Update(storage.Many(parsingIDs), storage.SetStatus(jobs.StatusQueued))
	if err != nil {
		return fmt.Errorf("could not update parsing jobs status: %w", err)
	}

	for _, parsing := range parsings {

		log := j.log.With().
			Uint64("chain_id", parsing.ChainID).
			Strs("contract_addresses", parsing.ContractAddresses).
			Strs("event_hashes", parsing.EventHashes).
			Uint64("start_height", parsing.StartHeight).
			Uint64("end_height", parsing.EndHeight).
			Str("job_status", string(parsing.JobStatus)).
			Logger()

		err := j.creator.PublishParsingJob(parsing)
		if err != nil {
			return fmt.Errorf("could not publish parsing job: %w", err)
		}

		log.Info().Msg("parsing job published")
	}

	return nil
}

func (j *Job) handleActionJobs(actions []*jobs.Action) error {

	actionIDs := make([]string, 0, len(actions))
	for _, action := range actions {
		actionIDs = append(actionIDs, action.ID)
	}

	err := j.actions.Update(storage.Many(actionIDs), storage.SetStatus(jobs.StatusQueued))
	if err != nil {
		return fmt.Errorf("could not update action jobs status: %w", err)
	}

	for _, action := range actions {

		log := j.log.With().
			Uint64("chain_id", action.ChainID).
			Str("contract_address", action.ContractAddress).
			Str("token_id", action.TokenID).
			Str("action_type", action.ActionType).
			Uint64("block_height", action.BlockHeight).
			Str("job_status", action.JobStatus).
			Logger()

		err := j.creator.PublishActionJob(action)
		if err != nil {
			return fmt.Errorf("could not publish action job: %w", err)
		}

		log.Info().Msg("action job published")
	}

	return nil
}
