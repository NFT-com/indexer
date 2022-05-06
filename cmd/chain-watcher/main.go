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

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/service/notifier"
	"github.com/NFT-com/indexer/service/notifier/heads"
	"github.com/NFT-com/indexer/service/notifier/multiplex"
	"github.com/NFT-com/indexer/service/notifier/ticker"
	"github.com/NFT-com/indexer/service/persister"
	"github.com/NFT-com/indexer/storage/filters"
	"github.com/NFT-com/indexer/storage/graph"
	storage "github.com/NFT-com/indexer/storage/jobs"
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
		flagLogLevel string

		flagGraphDB string
		flagJobsDB  string
		flagNodeURL string

		flagOpenConnections uint
		flagIdleConnections uint
		flagJobInterval     time.Duration
		flagPendingLimit    uint
		flagHeightRange     uint
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagGraphDB, "graph-database", "d", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagNodeURL, "node-url", "n", "ws://127.0.0.1:8545", "URL for Ethereum JSON RPC API connection")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of open database connections")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle database connections")
	pflag.DurationVar(&flagJobInterval, "write-interval", time.Second, "interval between checks for job writing")
	pflag.UintVar(&flagPendingLimit, "pending-limit", 1000, "maximum number of pending jobs per combination")
	pflag.UintVar(&flagHeightRange, "height-range", 10, "maximum heights to include in a single job")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	graphDB, err := sql.Open(params.DialectPostgres, flagGraphDB)
	if err != nil {
		return fmt.Errorf("could not open data DB: %w", err)
	}
	graphDB.SetMaxOpenConns(int(flagOpenConnections))
	graphDB.SetMaxIdleConns(int(flagIdleConnections))

	collectionRepo := graph.NewCollectionRepository(graphDB)
	standardRepo := graph.NewStandardRepository(graphDB)
	eventTypeRepo := graph.NewEventTypeRepository(graphDB)

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		return fmt.Errorf("could not open jobs DB: %w", err)
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := storage.NewParsingRepository(jobsDB)

	// Initialize the Ethereum node client and get the latest height to initialize
	// the watchers properly.
	client, err := ethclient.DialContext(ctx, flagNodeURL)
	if err != nil {
		return fmt.Errorf("could not connect to node: %w", err)
	}
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("could not get latest block number: %w", err)
	}

	collections, err := collectionRepo.Find()
	if err != nil {
		return fmt.Errorf("could not get collections: %w", err)
	}

	persist := persister.New(log, ctx, parsingRepo, 100*time.Millisecond, 1000)

	// For every collection and event type combination, initialize a ticker notifier
	// that will notify regularly of the latest height.
	var listens []notifier.Listener
	for _, collection := range collections {

		standards, err := standardRepo.Find(filters.Eq("collection_id", collection.ID))
		if err != nil {
			return fmt.Errorf("could not get standards (collection: %s)", collection.ID)
		}

		for _, standard := range standards {

			events, err := eventTypeRepo.Find(filters.Eq("standard_id", standard.ID))
			if err != nil {
				return fmt.Errorf("could not get event types: %w", err)
			}

			for _, event := range events {

				inputs := inputs.Parsing{
					NodeURL: flagNodeURL,
				}

				// create the job template for this combination
				template := jobs.Parsing{
					ID:          "",
					ChainID:     flagChainID,
					BlockNumber: 0,
					Address:     collection.Address,
					Standard:    standard.Name,
					Event:       event.ID,
					Status:      jobs.StatusCreated,
				}

				// initialize a job creator that will be notified of heights and
				// create jobs accordingly
				create := job.NewCreator(log, flagStartHeight, template, persist, jobsStore, flagPendingLimit)
				listens = append(listens, create)

				// initialize a ticker that will notify of the latest height at
				// regular intervals, to stay live when no blocks happen
				live, err := ticker.NewNotifier(log, ctx, flagJobInterval, latest, create)
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

	network, err := web3.New(ctx, flagNodeURL)
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
