package main

import (
	"database/sql"
	"os"
	"os/signal"
	"time"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/storage/graph"
	"github.com/NFT-com/indexer/storage/jobs"
	_ "github.com/lib/pq"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagLogLevel string

		flagGraphDB   string
		flagJobsDB    string
		flagRedisDB   int
		flagRedisURL  string
		flagAWSRegion string

		flagOpenConnections   uint
		flagIdleConnections   uint
		flagLambdaName        string
		flagLambdaConcurrency uint
		flagRateLimit         uint

		flagDryRun bool
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagGraphDB, "graph-database", "d", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagRedisURL, "redis-url", "u", "127.0.0.1:6379", "URL for Redis server connection")
	pflag.IntVarP(&flagRedisDB, "redis-database", "d", 1, "Redis database number")
	pflag.StringVarP(&flagAWSRegion, "aws-region", "r", "eu-west-1", "AWS region for Lambda invocation")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of open database connections")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle database connections")
	pflag.StringVar(&flagLambdaName, "lambda-name", "action-worker", "name of Lambda function for invocation")
	pflag.UintVar(&flagLambdaConcurrency, "lambda-concurrency", 900, "maximum number of concurrent Lambda invocations")
	pflag.UintVar(&flagRateLimit, "rate-limit", 100, "maximum number of API requests to Ethereum JSON RPC API")

	pflag.BoolVar(&flagDryRun, "dry-run", false, "executing as dry run disables invocation of Lambda function")

	pflag.Parse()

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("log_level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	sessionConfig := aws.Config{Region: aws.String(flagAWSRegion)}
	session := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(session)

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Str("jobs_database", flagJobsDB).Msg("could not connect to jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	actionRepo := jobs.NewActionRepository(jobsDB)

	graphDB, err := sql.Open(params.DialectPostgres, flagGraphDB)
	if err != nil {
		log.Error().Err(err).Str("graph_database", flagGraphDB).Msg("could not connect to graph database")
		return failure
	}
	graphDB.SetMaxOpenConns(int(flagOpenConnections))
	graphDB.SetMaxIdleConns(int(flagIdleConnections))

	collectionRepo := graph.NewCollectionRepository(graphDB)
	nftRepo := graph.NewNFTRepository(graphDB)
	traitRepo := graph.NewTraitRepository(graphDB)

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

	queue, err := rmqConnection.OpenQueue(params.QueueAction)
	if err != nil {
		log.Error().Err(err).Msg("could not open queue")
		return failure
	}

	err = queue.StartConsuming(int64(flagLambdaConcurrency), 200*time.Millisecond)
	if err != nil {
		log.Error().Err(err).Msg("could not start consuming")
		return failure
	}

	for i := uint(0); i < flagLambdaConcurrency; i++ {
		consumer := pipeline.NewActionConsumer(log, lambdaClient, flagLambdaName, actionRepo, collectionRepo, nftRepo, traitRepo, flagRateLimit, flagDryRun)
		_, err = queue.AddConsumer("action-consumer", consumer)
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
