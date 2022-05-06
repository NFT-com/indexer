package main

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/creator"
	"github.com/NFT-com/indexer/service/notifier"
	"github.com/NFT-com/indexer/storage/graph"
	storage "github.com/NFT-com/indexer/storage/jobs"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

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
		flagWriteInterval   time.Duration
		flagPendingLimit    uint
		flagHeightRange     uint
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagNodeURL, "node-url", "n", "ws://127.0.0.1:8545", "URL for Ethereum JSON RPC API connection")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 16, "maximum number of open database connections")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 4, "maximum number of idle database connections")
	pflag.DurationVar(&flagWriteInterval, "write-interval", time.Second, "interval between checks for job writing")
	pflag.UintVar(&flagPendingLimit, "pending-limit", 1000, "maximum number of pending jobs per combination")
	pflag.UintVar(&flagHeightRange, "height-range", 10, "maximum heights to include in a single job")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("log_level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	graphDB, err := sql.Open(params.DialectPostgres, flagGraphDB)
	if err != nil {
		log.Error().Err(err).Str("graph_db", flagGraphDB).Msg("could not open graph database")
		return failure
	}
	graphDB.SetMaxOpenConns(int(flagOpenConnections))
	graphDB.SetMaxIdleConns(int(flagIdleConnections))

	networkRepo := graph.NewNetworkRepository(graphDB)
	collectionRepo := graph.NewCollectionRepository(graphDB)

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Str("jobs_db", flagGraphDB).Msg("could not open jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := storage.NewParsingRepository(jobsDB)

	// Initialize the Ethereum node client and get the latest height to initialize
	// the watchers properly.
	client, err := ethclient.DialContext(ctx, flagNodeURL)
	if err != nil {
		log.Error().Err(err).Str("node_url", flagNodeURL).Msg("could not connect to node API")
		return failure
	}

	// Get all of the chain IDs from the graph database and initialize one creator
	// for each of the networks.
	networks, err := networkRepo.List()
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve list of networks")
		return failure
	}
	creators := make([]notifier.Listener, 0, len(networks))
	for _, network := range networks {
		creator := creator.New(log, collectionRepo, parsingRepo,
			creator.WithNodeURL(flagNodeURL),
			creator.WithChainID(network.ChainID),
			creator.WithPendingLimit(flagPendingLimit),
			creator.WithHeightRange(flagHeightRange),
		)
		creators = append(creators, creator)

		log.Info().
			Str("network", network.Name).
			Uint64("chain_id", network.ChainID).
			Msg("launching job creator")
	}

	// Initialize a multiplex notifier that will notify all of our creators at the
	// same time, a ticker notifier that will trigger it each interval, and a heads
	// notifier that will update its height.
	multi := notifier.NewMultiNotifier(creators...)
	ticker := notifier.NewTickerNotifier(log, ctx, multi,
		notifier.WithNotifyInterval(flagWriteInterval),
	)
	_, err = notifier.NewBlocksNotifier(log, ctx, client, ticker)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize blocks notifier")
		return failure
	}

	// Manually set the interval notifier to the latest height.
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not get latest block")
		return failure
	}
	ticker.Notify(latest)

	select {

	case <-ctx.Done():

	case <-sig:
		cancel()
	}

	return success
}
