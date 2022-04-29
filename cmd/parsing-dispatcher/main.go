package main

import (
	"database/sql"
	"fmt"
	"log"
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

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/queue/consumer/parsing"
	"github.com/NFT-com/indexer/service/postgres"
)

const (
	databaseDriver          = "postgres"
	defaultParsingQueueName = "parsing"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagBatchSize       int64
		flagJobsDB          string
		flagDataDB          string
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
		flagOpenConnections uint
		flagIdleConnections uint
	)

	// TODO: remove batch size and instead use time-based dispatching
	pflag.Int64VarP(&flagBatchSize, "batch", "b", 500, "batch size to process")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "", "jobs database connection string")
	pflag.StringVarP(&flagDataDB, "data-database", "d", "", "data database connection string")
	pflag.StringVarP(&flagEventsDB, "events-database", "e", "", "events database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVarP(&flagParsingQueue, "parsing-queue", "q", defaultParsingQueueName, "name of the queue for parsing")
	pflag.IntVarP(&flagRateLimit, "rate-limit", "t", 100, "maximum amount of lambdas that can be invoked per second")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database number")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "aws lambda region")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "parsing-agent", "rmq consumer tag")
	pflag.BoolVar(&flagDryRun, "dry-run", false, "when in dry run mode, no lambdas are invoked")
	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	sessionConfig := aws.Config{Region: aws.String(flagRegion)}
	session := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(session)

	dispatcher, err := function.New(lambdaClient)
	if err != nil {
		return fmt.Errorf("could not create function client dispatcher: %w", err)
	}

	jobDB, err := sql.Open(databaseDriver, flagJobsDB)
	if err != nil {
		return fmt.Errorf("could not open jobs SQL connection: %w", err)
	}
	jobDB.SetMaxOpenConns(int(flagOpenConnections))
	jobDB.SetMaxIdleConns(int(flagIdleConnections))

	jobStore, err := postgres.NewStore(jobDB)
	if err != nil {
		return fmt.Errorf("could not create job store: %w", err)
	}

	eventDB, err := sql.Open(databaseDriver, flagEventsDB)
	if err != nil {
		return fmt.Errorf("could not open events SQL connection: %w", err)
	}
	eventDB.SetMaxOpenConns(int(flagOpenConnections))
	eventDB.SetMaxIdleConns(int(flagIdleConnections))

	eventStore, err := postgres.NewStore(eventDB)
	if err != nil {
		return fmt.Errorf("could not create event store: %w", err)
	}

	dataDB, err := sql.Open(databaseDriver, flagDataDB)
	if err != nil {
		return fmt.Errorf("could not open data SQL connection: %w", err)
	}
	dataDB.SetMaxOpenConns(int(flagOpenConnections))
	dataDB.SetMaxIdleConns(int(flagIdleConnections))

	dataStore, err := postgres.NewStore(dataDB)
	if err != nil {
		return fmt.Errorf("could not create data store: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Network: flagRedisNetwork,
		Addr:    flagRedisURL,
		DB:      flagRedisDatabase,
	})

	failed := make(chan error)
	rmqConnection, err := rmq.OpenConnectionWithRedisClient(flagRMQTag, redisClient, failed)
	if err != nil {
		return fmt.Errorf("could not open redis connection: %w", err)
	}

	queue, err := rmqConnection.OpenQueue(flagParsingQueue)
	if err != nil {
		return fmt.Errorf("could not open redis queue: %w", err)
	}

	err = queue.StartConsuming(2*flagBatchSize, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("could not start consuming process: %w", err)
	}

	consumer := parsing.NewConsumer(log, dispatcher, jobStore, eventStore, dataStore, flagRateLimit, flagDryRun)
	consumerName, err := queue.AddBatchConsumer(flagRMQTag, flagBatchSize, 100*time.Millisecond, consumer)
	if err != nil {
		return fmt.Errorf("could not add parsing consumer: %w", err)
	}
	log = log.With().Str("name", consumerName).Logger()

	log.Info().Msg("started parsing dispatcher")
	consumer.Run()

	select {
	case <-sig:
		rmqConnection.StopAllConsuming()
		consumer.Close()
	case err := <-failed:
		return err
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced interruption")
	}()

	log.Info().Msg("stopped parsing dispatcher gracefully")

	return nil
}
