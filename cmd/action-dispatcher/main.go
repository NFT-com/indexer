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

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagActionQueue     string
		flagJobsDB          string
		flagGraphDB         string
		flagLogLevel        string
		flagRedisDatabase   int
		flagRedisNetwork    string
		flagRedisURL        string
		flagRegion          string
		flagRMQTag          string
		flagConsumerCount   uint
		flagRateLimit       uint
		flagOpenConnections uint
		flagIdleConnections uint
		flagDryRun          bool
	)

	pflag.StringVarP(&flagActionQueue, "action-queue", "q", params.QueueAction, "name of the queue for action jobs")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "", "jobs database connection string")
	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "", "data database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.IntVar(&flagRedisDatabase, "redis-database", 1, "redis database number")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.StringVar(&flagRegion, "aws-region", "eu-west-1", "aws lambda region")
	pflag.UintVar(&flagConsumerCount, "consumer-count", 900, "number of concurrent consumers for the parsing queue")
	pflag.UintVarP(&flagRateLimit, "rate-limit", "r", 100, "number of requests to the node per second")
	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")
	pflag.BoolVar(&flagDryRun, "dry-run", false, "whether to execute a dry run (don't invoke lambda)")

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

	sessionConfig := aws.Config{Region: aws.String(flagRegion)}
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

	queue, err := rmqConnection.OpenQueue(flagActionQueue)
	if err != nil {
		log.Error().Err(err).Msg("could not open queue")
		return failure
	}

	err = queue.StartConsuming(int64(flagConsumerCount), 200*time.Millisecond)
	if err != nil {
		log.Error().Err(err).Msg("could not start consuming")
		return failure
	}

	for i := uint(0); i < flagConsumerCount; i++ {
		consumer := pipeline.NewActionConsumer(log, lambdaClient, actionRepo, collectionRepo, nftRepo, traitRepo, flagRateLimit, flagDryRun)
		_, err = queue.AddConsumer("action_consumer", consumer)
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
