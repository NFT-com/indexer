package main

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"go.uber.org/ratelimit"

	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/nsqlog"
	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/storage/db"
	"github.com/NFT-com/indexer/storage/events"
	"github.com/NFT-com/indexer/storage/graph"
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
		flagLogLevel string

		flagGraphDB    string
		flagJobsDB     string
		flagEventsDB   string
		flagNSQLookups []string
		flagNSQServer  string
		flagLambdaName string

		flagOpenConnections   uint
		flagIdleConnections   uint
		flagHeightRange       uint
		flagRateLimit         uint
		flagLambdaConcurrency uint

		flagMinBackoff  time.Duration
		flagMaxBackoff  time.Duration
		flagMaxAttempts uint16

		flagDryRun bool
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")

	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.StringVarP(&flagEventsDB, "events-database", "e", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=events sslmode=disable", "Postgres connection details for events database")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagLambdaName, "lambda-name", "n", "parsing-worker", "name of the Lambda function to invoke")
	pflag.StringSliceVarP(&flagNSQLookups, "nsq-lookups", "k", []string{"127.0.0.1:4161"}, "address for NSQ lookup server to bootstrap consuming")
	pflag.StringVarP(&flagNSQServer, "nsq-server", "q", "127.0.0.1:4150", "address for NSQ server to produce messages")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")
	pflag.UintVar(&flagHeightRange, "height-range", 10, "maximum heights per parsing job")
	pflag.UintVar(&flagRateLimit, "rate-limit", 10, "maximum number of API requests per second")
	pflag.UintVar(&flagLambdaConcurrency, "lambda-concurrency", 100, "maximum number of concurrent Lambda invocations")

	pflag.DurationVar(&flagMinBackoff, "min-backoff", 1*time.Second, "minimum backoff duration for NSQ consumers")
	pflag.DurationVar(&flagMaxBackoff, "max-backoff", 10*time.Minute, "maximum backoff duration for NSQ consumers")
	pflag.Uint16Var(&flagMaxAttempts, "max-attempts", 3, "maximum number of attempts per job")

	pflag.BoolVar(&flagDryRun, "dry-run", false, "executing as dry run disables invocation of Lambda function")

	pflag.Parse()

	rand.Seed(time.Now().UnixNano())

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("log_level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("could not load AWS configuration")
		return failure
	}

	retrier := db.NewRetrier()

	graphDB, err := sql.Open(params.DialectPostgres, flagGraphDB)
	if err != nil {
		log.Error().Err(err).Str("graph_database", flagGraphDB).Msg("could not open graph database")
		return failure
	}
	graphDB.SetMaxOpenConns(int(flagOpenConnections))
	graphDB.SetMaxIdleConns(int(flagIdleConnections))

	collectionRepo := graph.NewCollectionRepository(graphDB)
	nftRepo := graph.NewNFTRepository(graphDB, retrier)
	ownerRepo := graph.NewOwnerRepository(graphDB, retrier)

	eventsDB, err := sql.Open(params.DialectPostgres, flagEventsDB)
	if err != nil {
		log.Error().Err(err).Str("events_database", flagEventsDB).Msg("could not connect to events database")
		return failure
	}
	eventsDB.SetMaxOpenConns(int(flagOpenConnections))
	eventsDB.SetMaxIdleConns(int(flagIdleConnections))

	transferRepo := events.NewTransferRepository(eventsDB, retrier)
	saleRepo := events.NewSaleRepository(eventsDB, retrier)

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Str("jobs_database", flagJobsDB).Msg("could not open jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	failureRepo := jobs.NewFailureRepository(jobsDB)

	nsqCfg := nsq.NewConfig()
	nsqCfg.MaxInFlight = int(flagLambdaConcurrency)
	nsqCfg.MaxAttempts = flagMaxAttempts
	nsqCfg.BackoffMultiplier = flagMinBackoff
	nsqCfg.MaxBackoffDuration = flagMaxBackoff
	consumer, err := nsq.NewConsumer(params.TopicParsing, params.ChannelDispatch, nsqCfg)
	if err != nil {
		log.Error().Err(err).Str("topic", params.TopicParsing).Str("channel", params.ChannelDispatch).Msg("could not create NSQ consumer")
		return failure
	}
	defer consumer.Stop()
	consumer.SetLogger(nsqlog.WrapForNSQ(log), nsq.LogLevelDebug)

	producer, err := nsq.NewProducer(flagNSQServer, nsqCfg)
	if err != nil {
		log.Error().Err(err).Str("nsq_address", flagNSQServer).Msg("could not create NSQ producer")
		return failure
	}
	defer producer.Stop()
	producer.SetLogger(nsqlog.WrapForNSQ(log), nsq.LogLevelDebug)

	lambda := lambda.NewFromConfig(awsCfg)
	limit := ratelimit.New(int(flagRateLimit))
	stage := pipeline.NewParsingStage(context.Background(), log, lambda, flagLambdaName, transferRepo, saleRepo, collectionRepo, nftRepo, ownerRepo, failureRepo, producer, limit,
		pipeline.WithParsingDryRun(flagDryRun),
		pipeline.WithParsingMaxAttempts(flagMaxAttempts),
	)
	consumer.AddConcurrentHandlers(stage, int(flagLambdaConcurrency))

	err = consumer.ConnectToNSQLookupds(flagNSQLookups)
	if err != nil {
		log.Error().Err(err).Strs("nsq_lookups", flagNSQLookups).Msg("could not connect to NSQ lookups")
		return failure
	}

	log.Info().Msg("parsing dispatcher started")

	<-sig

	log.Info().Msg("initialized shutdown")

	go func() {
		<-sig
		log.Fatal().Msg("forced shutdown")
	}()

	log.Info().Msg("shutdown complete")

	return success
}
