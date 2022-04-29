package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
	"github.com/NFT-com/indexer/queue/consumer/addition"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/service/postgres"
)

const (
	databaseDriver = "postgres"

	defaultHTTPTimeout       = time.Second * 30
	defaultPollDuration      = time.Second * 20
	defaultAdditionQueueName = "addition"
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
		flagAdditionQueueName    string
		flagAPIEndpoint          string
		flagConcurrentJobs       int
		flagConsumerPollDuration time.Duration
		flagConsumerPrefetch     int64
		flagDBConnectionInfo     string
		flagLogLevel             string
		flagRMQTag               string
		flagRedisDatabase        int
		flagRedisNetwork         string
		flagRedisURL             string
		flagRegion               string
	)

	pflag.StringVarP(&flagAdditionQueueName, "addition-queue", "q", defaultAdditionQueueName, "name of the queue for addition jobs")
	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base hostname and port")
	pflag.IntVarP(&flagConcurrentJobs, "jobs", "j", 4, "number of concurrent jobs for the consumer")
	pflag.DurationVarP(&flagConsumerPollDuration, "poll-duration", "i", defaultPollDuration, "consumer poll duration")
	pflag.Int64VarP(&flagConsumerPrefetch, "prefetch", "p", 5, "amount of messages to prefetch in the consumer")
	pflag.StringVarP(&flagDBConnectionInfo, "database", "d", "", "data database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "dispatcher-agent", "rmq consumer tag")
	pflag.IntVar(&flagRedisDatabase, "redis-database", 1, "redis database number")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "aws lambda region")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	db, err := sql.Open(databaseDriver, flagDBConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open data SQL connection: %w", err)
	}

	dataStore, err := postgres.NewStore(db)
	if err != nil {
		return fmt.Errorf("could not create data store: %w", err)
	}

	sessionConfig := aws.Config{Region: aws.String(flagRegion)}
	session := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(session)

	dispatcher, err := function.New(log, lambdaClient)
	if err != nil {
		return fmt.Errorf("could not create function client dispatcher: %w", err)
	}

	cli := http.DefaultClient
	cli.Timeout = defaultHTTPTimeout

	api := client.New(log,
		client.WithHTTPClient(cli),
		client.WithHost(flagAPIEndpoint),
	)

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

	queue, err := rmqConnection.OpenQueue(flagAdditionQueueName)
	if err != nil {
		return fmt.Errorf("could not open redis queue: %w", err)
	}

	err = queue.StartConsuming(flagConsumerPrefetch, flagConsumerPollDuration)
	if err != nil {
		return fmt.Errorf("could not start consuming process: %w", err)
	}

	consumer := addition.NewConsumer(log, api, dispatcher, dataStore, flagConcurrentJobs)
	consumerName, err := queue.AddConsumer(flagRMQTag, consumer)
	if err != nil {
		return fmt.Errorf("could not add addition consumer: %w", err)
	}
	log = log.With().Str("name", consumerName).Logger()

	log.Info().Msg("started addition dispatcher")
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

	log.Info().Msg("stopped addition dispatcher gracefully")

	return nil
}
