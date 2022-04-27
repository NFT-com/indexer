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
	"github.com/NFT-com/indexer/queue/consumer/action"
	"github.com/NFT-com/indexer/service/postgres"
)

const (
	databaseDriver         = "postgres"
	defaultActionQueueName = "action"
)

func main() {
	err := run()
	if err != nil {
		// TODO: Improve this mixing logging
		// https://github.com/NFT-com/indexer/issues/32
		log.Fatalln(err)
	}
}

func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagActionQueueName      string
		flagBatchSize            int64
		flagDBJobsConnectionInfo string
		flagDBDataConnectionInfo string
		flagLogLevel             string
		flagRateLimit            int
		flagRedisDatabase        int
		flagRedisNetwork         string
		flagRedisURL             string
		flagRegion               string
		flagRMQTag               string
	)

	pflag.StringVarP(&flagActionQueueName, "action-queue", "q", defaultActionQueueName, "name of the queue for action jobs")
	pflag.Int64VarP(&flagBatchSize, "batch", "b", 500, "batch size to process")
	pflag.StringVarP(&flagDBJobsConnectionInfo, "job-database", "j", "", "jobs database connection string")
	pflag.StringVarP(&flagDBDataConnectionInfo, "data-database", "d", "", "data database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.IntVar(&flagRedisDatabase, "redis-database", 1, "redis database number")
	pflag.IntVarP(&flagRateLimit, "rate-limit", "t", 100, "maximum amount of lambdas that can be invoked per second")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "aws lambda region")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "dispatcher-agent", "rmq consumer tag")

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

	jobDB, err := sql.Open(databaseDriver, flagDBJobsConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open jobs SQL connection: %w", err)
	}

	jobStore, err := postgres.NewStore(jobDB)
	if err != nil {
		return fmt.Errorf("could not create job store: %w", err)
	}

	dataDB, err := sql.Open(databaseDriver, flagDBDataConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open data SQL connection: %w", err)
	}

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

	queue, err := rmqConnection.OpenQueue(flagActionQueueName)
	if err != nil {
		return fmt.Errorf("could not open redis queue: %w", err)
	}

	err = queue.StartConsuming(2*flagBatchSize, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("could not start consuming process: %w", err)
	}

	consumer := action.NewConsumer(log, dispatcher, jobStore, dataStore, flagRateLimit)
	consumerName, err := queue.AddConsumer(flagRMQTag, consumer)
	if err != nil {
		return fmt.Errorf("could not add action consumer: %w", err)
	}
	log = log.With().Str("name", consumerName).Logger()

	log.Info().Msg("started action dispatcher")
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

	log.Info().Msg("stopped action dispatcher gracefully")

	return nil
}
