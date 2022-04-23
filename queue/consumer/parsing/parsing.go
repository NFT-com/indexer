package parsing

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/function/handlers/parsing"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/log"
	"github.com/NFT-com/indexer/service/client"
)

const concurrentConsumers = 4096

type Parsing struct {
	log           zerolog.Logger
	dispatcher    function.Invoker
	apiClient     *client.Client
	eventStore    Store
	dataStore     Store
	rateLimiter   <-chan time.Time
	consumerQueue chan [][]byte
	close         chan struct{}
}

func NewConsumer(log zerolog.Logger, apiClient *client.Client, dispatcher function.Invoker, eventStore Store, dataStore Store, rateLimit int) *Parsing {
	limiter := time.Tick(time.Second / time.Duration(rateLimit))

	c := Parsing{
		log:           log,
		dispatcher:    dispatcher,
		apiClient:     apiClient,
		eventStore:    eventStore,
		dataStore:     dataStore,
		rateLimiter:   limiter,
		consumerQueue: make(chan [][]byte, concurrentConsumers),
		close:         make(chan struct{}),
	}

	return &c
}

func (d *Parsing) Consume(batch rmq.Deliveries) {
	payloads := make([][]byte, 0, len(batch))

	for _, delivery := range batch {
		payload := []byte(delivery.Payload())
		payloads = append(payloads, payload)

		err := delivery.Ack()
		if err != nil {
			d.log.Error().Err(err).Msg("could not acknowledge message")
			return
		}
	}

	d.consumerQueue <- payloads
}

func (d *Parsing) Run() {
	for i := 0; i < concurrentConsumers; i++ {
		go func() {
			for {
				select {
				case <-d.close:
					return
				case payload := <-d.consumerQueue:
					d.consume(payload)
				}
			}
		}()
	}
}

func (d *Parsing) Close() {
	close(d.close)
}

func (d *Parsing) consume(payloads [][]byte) {
	jobList := d.unmarshalJobs(payloads)

	inputMap := make(map[uint64]parsing.Input, len(jobList))

	// Job List is unordered, in order to group them in a continuous way,
	// first the next loop basically maps the height to an input that could group jobs
	// from the same height. It also gets the highest and lowest heights in the list.
	lowestBlock := uint64(math.MaxUint64)
	highestBlock := uint64(0)
	for _, job := range jobList {
		// gets the current block height
		currentBlock, err := blockToUint64(job.BlockNumber)
		if err != nil {
			d.handleError(err, "could not parse block into uint64", job.ID)
			continue
		}

		// checks if there is already an entry, if so joins the ids, addresses, event_types and maps the standard.
		input, ok := inputMap[currentBlock]
		if ok {
			input.IDs = append(input.IDs, job.ID)
			input.Addresses = append(input.Addresses, job.Address)
			input.EventTypes = append(input.EventTypes, job.EventType)
			input.Standards[job.EventType] = job.StandardType

			inputMap[currentBlock] = input
			continue
		}

		// if there is no entry just create a new one and check if this is lower that
		// the current lowest height or upper than the highest height.
		inputMap[currentBlock] = parsing.Input{
			IDs:        []string{job.ID},
			ChainURL:   job.ChainURL,
			ChainID:    job.ChainID,
			ChainType:  job.ChainType,
			StartBlock: job.BlockNumber,
			EndBlock:   job.BlockNumber,
			Addresses:  []string{job.Address},
			Standards:  map[string]string{job.EventType: job.StandardType},
			EventTypes: []string{job.EventType},
		}

		if lowestBlock > currentBlock {
			lowestBlock = currentBlock
		}

		if highestBlock < currentBlock {
			highestBlock = currentBlock
		}
	}

	inputs := make([]parsing.Input, 0)

	// Merges all the continuous inputs heights into a single one.
	input := parsing.Input{}
	for i := lowestBlock; i <= highestBlock; i++ {
		part, ok := inputMap[i]
		if !ok {
			if len(input.IDs) != 0 {
				inputs = append(inputs, input)
			}

			input = parsing.Input{}
			continue
		}

		// If current input ids len is 0, it leans it does not have a current input
		if len(input.IDs) == 0 {
			input = part
			continue
		}

		input = d.fillInput(input, part)
	}

	if len(input.IDs) != 0 {
		inputs = append(inputs, input)
	}

	for _, input := range inputs {
		payload, err := json.Marshal(input)
		if err != nil {
			d.handleError(err, "could not marshal input payload", input.IDs...)
			return
		}

		// Wait for rate limiter to have available spots.
		<-d.rateLimiter

		name := functionName(input)

		notify := func(err error, dur time.Duration) {
			d.log.Error().
				Err(err).
				Dur("retry_in", dur).
				Str("name", name).
				Int("payload_len", len(payload)).
				Msg("count not invoke lambda")
		}
		var output []byte
		_ = backoff.RetryNotify(func() error {
			output, err = d.dispatcher.Invoke(name, payload)
			if err != nil {
				return err
			}
			return nil
		}, backoff.NewExponentialBackOff(), notify)

		var logs []log.Log
		err = json.Unmarshal(output, &logs)
		if err != nil {
			d.handleError(err, "could not unmarshal output logs", input.IDs...)
			return
		}

		err = d.processLogs(input, logs)
		if err != nil {
			d.handleError(err, "could not handle output logs", input.IDs...)
			return
		}

		for _, id := range input.IDs {
			err = d.apiClient.UpdateParsingJobStatus(id, jobs.StatusFinished)
			if err != nil {
				d.handleError(err, "could not update job status", id)
				continue
			}
		}
	}
}

func (d *Parsing) fillInput(base, part parsing.Input) parsing.Input {
	base.EndBlock = part.EndBlock
	base.IDs = append(base.IDs, part.IDs...)

	addresses := make(map[string]struct{}, len(base.Addresses))
	for _, address := range base.Addresses {
		addresses[address] = struct{}{}
	}

	for _, address := range part.Addresses {
		if _, ok := addresses[address]; ok {
			continue
		}

		base.Addresses = append(base.Addresses, address)
	}

	eventTypes := make(map[string]struct{}, len(base.EventTypes))
	for _, eventType := range base.EventTypes {
		eventTypes[eventType] = struct{}{}
	}

	for _, eventType := range part.EventTypes {
		if _, ok := eventTypes[eventType]; ok {
			continue
		}

		base.EventTypes = append(base.EventTypes, eventType)
	}

	return base
}

func (d *Parsing) unmarshalJobs(payloads [][]byte) []jobs.Parsing {
	jobList := make([]jobs.Parsing, 0, len(payloads))

	for _, payload := range payloads {
		var job jobs.Parsing
		err := json.Unmarshal(payload, &job)
		if err != nil {
			d.log.Error().Err(err).Msg("could not unmarshal message")
			continue
		}

		storedJob, err := d.apiClient.GetParsingJob(job.ID)
		if err != nil {
			d.handleError(err, "could not retrieve parsing job", job.ID)
			continue
		}

		// job has been canceled meanwhile, no need to go further
		if job.Status != jobs.StatusCreated && storedJob.Status != jobs.StatusCanceled {
			continue
		}

		err = d.apiClient.UpdateParsingJobStatus(job.ID, jobs.StatusProcessing)
		if err != nil {
			d.handleError(err, "could not update job status", job.ID)
			continue
		}

		jobList = append(jobList, job)
	}

	return jobList
}

func (d *Parsing) handleError(err error, message string, ids ...string) {
	for _, id := range ids {
		updateErr := d.apiClient.UpdateParsingJobStatus(id, jobs.StatusFailed)
		if updateErr != nil {
			d.log.Error().Err(updateErr).Msg("could not update job status")
		}
	}

	d.log.Error().Err(err).Strs("job_id", ids).Msg(message)
}

func blockToUint64(block string) (uint64, error) {
	number := big.NewInt(0)
	_, ok := number.SetString(block, 0)
	if !ok {
		return 0, fmt.Errorf("could not parse block into big.Int")
	}

	return number.Uint64(), nil
}

func functionName(input parsing.Input) string {
	h := sha256.New()

	s := strings.Join(
		[]string{
			"parsing",
			strings.ToLower(input.ChainType),
		},
		"-",
	)
	h.Write([]byte(s))

	name := fmt.Sprintf("%x", h.Sum(nil))

	return name
}
