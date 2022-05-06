package watcher

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/storage/statements"
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
		log:      log.With().Str("component", "watcher").Logger(),
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

			parsings, err := j.parsings.Find(statements.Eq("status", jobs.StatusCreated))
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve parsing jobs")
				continue
			}
			if len(parsings) == 0 {
				continue
			}

			j.handleParsingJobs(parsings)

		case <-j.close:
			return
		}
	}
}

func (j *Job) watchActions() {

	for {
		select {

		case <-time.After(j.delay):

			actions, err := j.actions.Find(statements.Eq("status", jobs.StatusCreated))
			if err != nil {
				j.log.Error().Err(err).Msg("could not retrieve action jobs")
				continue
			}
			if len(actions) == 0 {
				continue
			}

			j.handleActionJobs(actions)

		case <-j.close:
			return
		}
	}
}

func (j *Job) handleParsingJobs(parsings []*jobs.Parsing) {

	for _, parsing := range parsings {
		err := j.publishParsingJob(parsing)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("parsing_id", parsing.ID).
				Uint64("chain_id", parsing.ChainID).
				Strs("contract_addresses", parsing.ContractAddresses).
				Strs("event_hashes", parsing.EventHashes).
				Str("status", string(parsing.Status)).
				Msg("could not publish parsing job")
			continue
		}
	}

	j.log.Info().Int("parsings", len(parsings)).Msg("queued parsing jobs")
}

func (j *Job) publishParsingJob(parsing *jobs.Parsing) error {

	if parsing.Status != jobs.StatusCreated {
		return nil
	}

	err := j.produce.PublishParsingJob(parsing)
	if err != nil {
		return fmt.Errorf("could not get publish parsing job: %w", err)
	}

	err = j.parsings.UpdateStatus(parsing.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update parsing job status: %w", err)
	}

	return nil
}

func (j *Job) handleActionJobs(actions []*jobs.Action) {

	for _, action := range actions {
		err := j.publishActionJob(action)
		if err != nil {
			j.log.Error().
				Err(err).
				Str("action_id", action.ID).
				Uint64("chain_id", action.ChainID).
				Str("address", action.Address).
				Str("token_id", action.TokenID).
				Str("status", string(action.Status)).
				Msg("could not publish action job")
			continue
		}
	}

	j.log.Info().Int("actions", len(actions)).Msg("queued action jobs")
}

func (j *Job) publishActionJob(action *jobs.Action) error {

	if action.Status != jobs.StatusCreated {
		return nil
	}

	err := j.produce.PublishActionJob(action)
	if err != nil {
		return fmt.Errorf("could not publish action job: %w", err)
	}

	err = j.actions.UpdateStatus(action.ID, jobs.StatusQueued)
	if err != nil {
		return fmt.Errorf("could not update action job status: %w", err)
	}

	return nil
}
