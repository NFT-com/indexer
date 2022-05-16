package main

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

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

		flagGraphDB      string
		flagJobDB        string
		flagNodeURL      string
		flagWebsocketURL string

		flagOpenConnections uint
		flagIdleConnections uint
		flagWriteInterval   time.Duration
		flagPendingLimit    uint
		flagHeightRange     uint
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "postgresql connection details for graph database")
	pflag.StringVarP(&flagJobDB, "job-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "postgresql connection details for job database")
	pflag.StringVarP(&flagNodeURL, "node-url", "n", "http://127.0.0.1:8545", "http URL for Ethereum JSON RPC API connection")
	pflag.StringVarP(&flagWebsocketURL, "websocket-url", "w", "ws://127.0.0.1:8545", "websocket URL for Ethereum JSON RPC API connection")

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
	marketplaceRepo := graph.NewMarketplaceRepository(graphDB)

	jobDB, err := sql.Open(params.DialectPostgres, flagJobDB)
	if err != nil {
		log.Error().Err(err).Str("job_db", flagGraphDB).Msg("could not open job database")
		return failure
	}
	jobDB.SetMaxOpenConns(int(flagOpenConnections))
	jobDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := storage.NewParsingRepository(jobDB)

	// Get all of the chain IDs from the graph database and initialize one creator
	// for each of the networks.
	networks, err := networkRepo.List()
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve list of networks")
		return failure
	}
	creators := make([]notifier.Listener, 0, len(networks))
	for _, network := range networks {
		creator := creator.New(log, collectionRepo, marketplaceRepo, parsingRepo,
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
	_, err = notifier.NewBlocksNotifier(log, ctx, flagWebsocketURL, ticker)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize blocks notifier")
		return failure
	}

	// Initialize the Ethereum node client and get the latest height.
	client, err := ethclient.DialContext(ctx, flagNodeURL)
	if err != nil {
		log.Error().Err(err).Str("node_url", flagNodeURL).Msg("could not connect to node API")
		return failure
	}

	// Manually set the interval notifier to the latest height.
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not get latest block")
		return failure
	}
	ticker.Notify(latest)

	log.Info().Msg("job creator started")
	select {

	case <-ctx.Done():

	case <-sig:
		log.Info().Msg("job creator stopping")
		cancel()
	}

	log.Info().Msg("job creator done")

	return success
}
