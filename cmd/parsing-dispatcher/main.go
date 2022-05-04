package main

import (
	"database/sql"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/storage/events"
	"github.com/NFT-com/indexer/storage/jobs"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagJobsDB          string
		flagEventsDB        string
		flagLogLevel        string
		flagParsingQueue    string
		flagRateLimit       int
		flagRedisDatabase   int
		flagRedisNetwork    string
		flagRedisURL        string
		flagRegion          string
		flagRMQTag          string
		flagDryRun          bool
		flagConsumerCount   uint
		flagOpenConnections uint
		flagIdleConnections uint
	)

	// TODO: remove batch size and instead use time-based dispatching
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "", "jobs database connection string")
	pflag.StringVarP(&flagEventsDB, "events-database", "e", "", "events database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVarP(&flagParsingQueue, "parsing-queue", "q", params.QueueParsing, "name of the queue for parsing")
	pflag.IntVarP(&flagRateLimit, "rate-limit", "t", 100, "maximum amount of lambdas that can be invoked per second")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database number")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "aws lambda region")
	pflag.BoolVar(&flagDryRun, "dry-run", false, "when in dry run mode, no lambdas are invoked")
	pflag.UintVar(&flagConsumerCount, "consumer-count", 100, "number of concurrent consumers for the parsing queue")
	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	sessionConfig := aws.Config{Region: aws.String(flagRegion)}
	session := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(session)

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := jobs.NewParsingRepository(jobsDB)
	actionRepo := jobs.NewActionRepository(jobsDB)

	eventsDB, err := sql.Open(params.DialectPostgres, flagEventsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to events database")
		return failure
	}
	eventsDB.SetMaxOpenConns(int(flagOpenConnections))
	eventsDB.SetMaxIdleConns(int(flagIdleConnections))

	mintRepo := events.NewMintRepository(eventsDB)
	transferRepo := events.NewTransferRepository(eventsDB)
	saleRepo := events.NewSaleRepository(eventsDB)
	burnRepo := events.NewBurnRepository(eventsDB)

	redisClient := redis.NewClient(&redis.Options{
		Network: flagRedisNetwork,
		Addr:    flagRedisURL,
		DB:      flagRedisDatabase,
	})
	failed := make(chan error)
	rmqConnection, err := rmq.OpenConnectionWithRedisClient(flagRMQTag, redisClient, failed)
	if err != nil {
		log.Error().Err(err).Str("redis_url", flagRedisURL).Msg("could not connect to redis server")
		return failure
	}
	defer rmqConnection.StopAllConsuming()

	queue, err := rmqConnection.OpenQueue(flagParsingQueue)
	if err != nil {
		log.Error().Err(err).Msg("could not open queue")
		return failure
	}

	err = queue.StartConsuming(int64(flagConsumerCount), 200*time.Millisecond)
	if err != nil {
		log.Error().Err(err).Msg("could not start consuming queue")
		return failure
	}

	for i := uint(0); i < flagConsumerCount; i++ {
		consumer := pipeline.NewParsingConsumer(log, lambdaClient, parsingRepo, actionRepo, mintRepo, transferRepo, saleRepo, burnRepo, flagRateLimit, flagDryRun)
		_, err := queue.AddConsumer("parsing_consumer", consumer)
		if err != nil {
			log.Error().Err(err).Msg("could not add consumer")
			return failure
		}
	}

	select {
	case <-sig:
		log.Info().Msg("initialized shutdown")
	case err = <-failed:
		return failure
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced shutdown")
	}()

	log.Info().Msg("shutdown complete")

	return success
}
