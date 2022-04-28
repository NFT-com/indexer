package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/creator/job"
	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/notifier"
	"github.com/NFT-com/indexer/notifier/heads"
	"github.com/NFT-com/indexer/notifier/multiplex"
	"github.com/NFT-com/indexer/notifier/ticker"
	"github.com/NFT-com/indexer/persister/database"
	"github.com/NFT-com/indexer/service/postgres"
)

const (
	databaseDriver = "postgres"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Command line parameter initialization.
	var (
		flagChainID              string
		flagChainURL             string
		flagChainType            string
		flagDBDataConnectionInfo string
		flagDBJobsConnectionInfo string
		flagLogLevel             string
		flagStartHeight          uint64
		flagJobLimit             uint
		flagNotifyPeriod         time.Duration
	)

	pflag.StringVarP(&flagChainID, "chain-id", "i", "", "id of the chain")
	pflag.StringVarP(&flagChainURL, "chain-url", "u", "", "url of the chain to connect")
	pflag.StringVarP(&flagChainType, "chain-type", "t", "", "type of chain")
	pflag.StringVarP(&flagDBDataConnectionInfo, "data-database", "d", "", "data database connection string")
	pflag.StringVarP(&flagDBJobsConnectionInfo, "job-database", "j", "", "jobs database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Uint64VarP(&flagStartHeight, "start-height", "s", 0, "default start height when no jobs found")
	pflag.UintVarP(&flagJobLimit, "job-limit", "b", 1000, "maximum number of pending jobs per combination")
	pflag.DurationVarP(&flagNotifyPeriod, "notify-period", "c", time.Second, "how often to notify watchers to create new jobs")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	// Open the jobs SQL database connection.
	dataDB, err := sql.Open(databaseDriver, flagDBDataConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open SQL connection: %w", err)
	}

	// Initialize the persister that will store jobs to the DB.
	dataStore, err := postgres.NewStore(dataDB)

	// Initialize the Ethereum node client and get the latest height to initialize
	// the watchers properly.
	client, err := ethclient.DialContext(ctx, flagChainURL)
	if err != nil {
		return fmt.Errorf("could not connect to node: %w", err)
	}
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("could not get latest block number: %w", err)
	}

	// Get the collections for the configured chain ID from the database.
	chain, err := dataStore.Chain(flagChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}
	collections, err := dataStore.Collections(chain.ID)
	if err != nil {
		return fmt.Errorf("could not get collections from store (chain: %s): %w", flagChainID, err)
	}

	// Open the jobs SQL database connection.
	jobDB, err := sql.Open(databaseDriver, flagDBJobsConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open SQL connection: %w", err)
	}

	// Initialize the persister that will store jobs to the DB.
	jobStore, err := postgres.NewStore(jobDB)
	if err != nil {
		return fmt.Errorf("could not create store: %w", err)
	}
	persist := database.NewPersister(log, ctx, jobStore, 100*time.Millisecond, 1000)

	// For every collection and event type combination, initialize a ticker notifier
	// that will notify regularly of the latest height.
	var listens []notifier.Listener
	for _, collection := range collections {

		standards, err := dataStore.Standards(collection.ID)
		if err != nil {
			return fmt.Errorf("could not get standards (collection: %s)", collection.ID)
		}

		for _, standard := range standards {

			events, err := dataStore.EventTypes(standard.ID)
			if err != nil {
				return fmt.Errorf("could not get event types: %w", err)
			}

			for _, event := range events {

				// TODO: retrieve starting height from latest completed job table
				// TODO: fix block number to be integer and rest of parsing job fields

				// create the job template for this combination
				template := jobs.Parsing{
					ID:          "",
					ChainURL:    flagChainURL,
					ChainID:     flagChainID,
					ChainType:   flagChainType,
					BlockNumber: "",
					Address:     collection.Address,
					Standard:    standard.Name,
					Event:       event.ID,
					Status:      jobs.StatusCreated,
				}

				// initialize a job creator that will be notified of heights and
				// create jobs accordingly
				create := job.NewCreator(log, flagStartHeight, template, persist, jobStore, flagJobLimit)
				listens = append(listens, create)

				// TODO: introduce jitter so not all DB requests hit it at the same second

				// initialize a ticker that will notify of the latest height at
				// regular intervals, to stay live when no blocks happen
				live, err := ticker.NewNotifier(log, ctx, flagNotifyPeriod, latest, create)
				if err != nil {
					return fmt.Errorf("could not create live notifier: %w", err)
				}
				listens = append(listens, live)
			}
		}
	}
	multi := multiplex.NewNotifier(listens...)
	_, err = heads.NewNotifier(log, ctx, client, multi)
	if err != nil {
		return fmt.Errorf("could not initialize heads notifier: %w", err)
	}

	network, err := web3.New(ctx, flagChainURL)
	if err != nil {
		return fmt.Errorf("could not create web3 network: %w", err)
	}

	chainID, err := network.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("could not get chain ID from network: %w", err)
	}
	if chainID != flagChainID {
		return fmt.Errorf("could not start watcher: mismatch between chain ID and chain URL")
	}

	select {

	case <-ctx.Done():

	case <-sig:
		cancel()
	}

	return nil
}
