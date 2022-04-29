package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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

	// Seed random generator for jitter.
	rand.Seed(time.Now().UnixNano())

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Command line parameter initialization.
	var (
		flagChainID         string
		flagChainURL        string
		flagChainType       string
		flagDataDB          string
		flagJobsDB          string
		flagLogLevel        string
		flagStartHeight     uint64
		flagJobLimit        uint
		flagNotifyPeriod    time.Duration
		flagOpenConnections uint
		flagIdleConnections uint
	)

	pflag.StringVarP(&flagChainID, "chain-id", "i", "", "id of the chain")
	pflag.StringVarP(&flagChainURL, "chain-url", "u", "", "url of the chain to connect")
	pflag.StringVarP(&flagChainType, "chain-type", "t", "", "type of chain")
	pflag.StringVarP(&flagDataDB, "data-database", "d", "", "connection details for data DB")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "", "connection details for jobs DB")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Uint64VarP(&flagStartHeight, "start-height", "s", 0, "default start height when no jobs found")
	pflag.UintVar(&flagJobLimit, "job-limit", 1000, "maximum number of pending jobs per combination")
	pflag.DurationVarP(&flagNotifyPeriod, "notify-period", "c", 100*time.Millisecond, "how often to notify watchers to create new jobs")
	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 70, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 20, "maximum number of idle connections")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	dataDB, err := sql.Open(databaseDriver, flagDataDB)
	if err != nil {
		return fmt.Errorf("could not open data DB: %w", err)
	}
	dataDB.SetMaxOpenConns(int(flagOpenConnections))
	dataDB.SetMaxIdleConns(int(flagIdleConnections))

	jobsDB, err := sql.Open(databaseDriver, flagJobsDB)
	if err != nil {
		return fmt.Errorf("could not open jobs DB: %w", err)
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

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

	// Get the collections for the configured chain ID from the database and initialize
	// the persister that will store jobs to the DB.
	dataStore, err := postgres.NewStore(dataDB)
	if err != nil {
		return fmt.Errorf("could not create data store: %w", err)
	}
	chain, err := dataStore.Chain(flagChainID)
	if err != nil {
		return fmt.Errorf("could not get chain: %w", err)
	}
	collections, err := dataStore.Collections(chain.ID)
	if err != nil {
		return fmt.Errorf("could not get collections from store (chain: %s): %w", flagChainID, err)
	}

	jobsStore, err := postgres.NewStore(jobsDB)
	if err != nil {
		return fmt.Errorf("could not create jobs store: %w", err)
	}
	persist := database.NewPersister(log, ctx, jobsStore, 100*time.Millisecond, 1000)

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

				// TODO: fix block number to be integer

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
				create := job.NewCreator(log, flagStartHeight, template, persist, jobsStore, flagJobLimit)
				listens = append(listens, create)

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
