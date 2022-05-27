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
		flagLogLevel string

		flagJobsDB     string
		flagEventsDB   string
		flagRedisDB    int
		flagRedisURL   string
		flagAWSRegion  string
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
	pflag.StringVarP(&flagRedisURL, "redis-url", "u", "127.0.0.1:6379", "Redis server URL")
	pflag.IntVarP(&flagRedisDB, "redis-database", "d", 1, "Redis database number")
	pflag.StringVarP(&flagAWSRegion, "aws-region", "r", "eu-west-1", "AWS region for Lambda invocation")
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

	sessionConfig := aws.Config{Region: aws.String(flagAWSRegion)}
	session := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(session)

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to job database")
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

	transferRepo := events.NewTransferRepository(eventsDB)
	saleRepo := events.NewSaleRepository(eventsDB)

	redisClient := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    flagRedisURL,
		DB:      flagRedisDB,
	})
	failed := make(chan error)
	rmqConnection, err := rmq.OpenConnectionWithRedisClient(params.PipelineIndexer, redisClient, failed)
	if err != nil {
		log.Error().Err(err).Str("redis_url", flagRedisURL).Msg("could not connect to redis server")
		return failure
	}
	defer rmqConnection.StopAllConsuming()

	queue, err := rmqConnection.OpenQueue(params.QueueParsing)
	if err != nil {
		log.Error().Err(err).Msg("could not open queue")
		return failure
	}

	err = queue.StartConsuming(int64(flagLambdaConcurrency), 200*time.Millisecond)
	if err != nil {
		log.Error().Err(err).Msg("could not start consuming queue")
		return failure
	}

	for i := uint(0); i < flagLambdaConcurrency; i++ {
		consumer := pipeline.NewParsingConsumer(log, lambdaClient, flagLambdaName, parsingRepo, actionRepo, transferRepo, saleRepo, flagRateLimit, flagDryRun)
		_, err := queue.AddConsumer("parsing-consumer", consumer)
		if err != nil {
			log.Error().Err(err).Msg("could not add consumer")
			return failure
		}
	}

	select {
	case <-sig:
		log.Info().Msg("initialized shutdown")
	case err = <-failed:
		log.Error().Err(err).Msg("execution failed")
		return failure
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced shutdown")
	}()

	log.Info().Msg("shutdown complete")

	return success
}
