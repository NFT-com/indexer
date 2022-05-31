package main

import (
	"context"
	"database/sql"
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
		flagLogLevel string

		flagJobsDB     string
		flagEventsDB   string
		flagNSQLookups []string
		flagLambdaName string

		flagOpenConnections   uint
		flagIdleConnections   uint
		flagHeightRange       uint
		flagRateLimit         uint
		flagLambdaConcurrency uint

		flagDryRun bool
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")

	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagEventsDB, "events-database", "e", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=events sslmode=disable", "Postgres connection details for events database")
	pflag.StringSliceVarP(&flagNSQLookups, "nsq-lookups", "q", []string{"127.0.0.1:4161"}, "address for NSQ lookup server to bootstrap consuming")
	pflag.StringVarP(&flagLambdaName, "lambda-name", "n", "parsing-worker", "name of the Lambda function to invoke")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")
	pflag.UintVar(&flagHeightRange, "height-range", 10, "maximum heights per parsing job")
	pflag.UintVar(&flagRateLimit, "rate-limit", 10, "maximum number of API requests per second")
	pflag.UintVar(&flagLambdaConcurrency, "lambda-concurrency", 100, "maximum number of concurrent Lambda invocations")

	pflag.BoolVar(&flagDryRun, "dry-run", false, "executing as dry run disables invocation of Lambda function")

	pflag.Parse()

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("could not load AWS configuration")
		return failure
	}

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Str("jobs_database", flagJobsDB).Msg("could not connect to jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := jobs.NewParsingRepository(jobsDB)
	actionRepo := jobs.NewActionRepository(jobsDB)

	eventsDB, err := sql.Open(params.DialectPostgres, flagEventsDB)
	if err != nil {
		log.Error().Err(err).Str("events_database", flagEventsDB).Msg("could not connect to events database")
		return failure
	}
	eventsDB.SetMaxOpenConns(int(flagOpenConnections))
	eventsDB.SetMaxIdleConns(int(flagIdleConnections))

	transferRepo := events.NewTransferRepository(eventsDB)
	saleRepo := events.NewSaleRepository(eventsDB)

	config := nsq.NewConfig()
	err = config.Set("max-in-flight", flagLambdaConcurrency*10)
	if err != nil {
		log.Error().Err(err).Uint("max-in-flight", 2*flagLambdaConcurrency).Msg("could not update NSQ configuration")
		return failure
	}

	consumer, err := nsq.NewConsumer(params.TopicParsing, params.PipelineIndexer, config)
	if err != nil {
		log.Error().Err(err).Str("topic", params.TopicParsing).Str("channel", params.PipelineIndexer).Msg("could not create NSQ consumer")
		return failure
	}

	lambda := lambda.NewFromConfig(cfg)
	limit := ratelimit.New(int(flagRateLimit))
	handler := pipeline.NewParsingHandler(context.Background(), log, lambda, flagLambdaName, parsingRepo, actionRepo, transferRepo, saleRepo, limit, flagDryRun)
	consumer.AddConcurrentHandlers(handler, int(flagLambdaConcurrency))

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
