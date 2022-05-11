package watcher

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/service/pipeline"
)

type Job struct {
	log      zerolog.Logger
	parsings ParsingStore
	actions  ActionStore
	produce  *pipeline.Producer
	delay    time.Duration
	close    chan struct{}
}

func New(log zerolog.Logger, parsings ParsingStore, actions ActionStore, produce *pipeline.Producer, delay time.Duration) *Job {

	j := Job{
		log:      log.With().Str("component", "jobs_watcher").Logger(),
		parsings: parsings,
		actions:  actions,
		produce:  produce,
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

		log := j.log.With().
			Uint64("chain_id", parsing.ChainID).
			Strs("contract_addresses", parsing.ContractAddresses).
			Strs("event_hashes", parsing.EventHashes).
			Uint64("start_height", parsing.StartHeight).
			Uint64("end_height", parsing.EndHeight).
			Str("job_status", string(parsing.JobStatus)).
			Logger()

		err := j.produce.PublishParsingJob(parsing)
		if err != nil {
			return fmt.Errorf("could not publish parsing job: %w", err)
		}

		parsingIDs = append(parsingIDs, parsing.ID)

		log.Info().Msg("parsing job published")
	}

	err := j.parsings.UpdateStatus(jobs.StatusQueued, "", parsingIDs...)
	if err != nil {
		return fmt.Errorf("could not update parsing jobs status: %w", err)
	}

	return nil
}

func (j *Job) handleActionJobs(actions []*jobs.Action) error {

	actionIDs := make([]string, 0, len(actions))
	for _, action := range actions {

		log := j.log.With().
			Uint64("chain_id", action.ChainID).
			Str("contract_address", action.ContractAddress).
			Str("token_id", action.TokenID).
			Str("action_type", action.ActionType).
			Uint64("block_height", action.BlockHeight).
			Str("job_status", action.JobStatus).
			Logger()

		err := j.produce.PublishActionJob(action)
		if err != nil {
			return fmt.Errorf("could not publish action job: %w", err)
		}

		actionIDs = append(actionIDs, action.ID)

		log.Info().Msg("action job published")
	}

	err := j.actions.UpdateStatus(jobs.StatusQueued, "", actionIDs...)
	if err != nil {
		return fmt.Errorf("could not update action jobs status: %w", err)
	}

	return nil
}
