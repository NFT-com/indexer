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
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/queue/consumer"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/service/postgres"
)

const (
	databaseDriver = "postgres"

	defaultHTTPTimeout      = time.Second * 30
	defaultPollDuration     = time.Second * 20
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
		flagAPIEndpoint          string
		flagConsumerPrefetch     int64
		flagConsumerPollDuration time.Duration
		flagDBConnectionInfo     string
		flagParsingQueueName     string
		flagLogLevel             string
		flagRMQTag               string
		flagRedisDatabase        int
		flagRedisNetwork         string
		flagRedisURL             string
		flagRegion               string
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base hostname and port")
	pflag.Int64VarP(&flagConsumerPrefetch, "prefetch", "p", 5, "amount of message to prefetch in the consumer")
	pflag.DurationVarP(&flagConsumerPollDuration, "poll-duration", "i", defaultPollDuration, "time for each consumer poll")
	pflag.StringVarP(&flagDBConnectionInfo, "db", "d", "", "database connection string")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "q", defaultParsingQueueName, "name of the queue for parsing")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "parsing-agent", "rmq consumer tag")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database number")
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

	sessionConfig := aws.Config{Region: aws.String(flagRegion)}
	session := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(session)

	dispatcher, err := function.New(lambdaClient)
	if err != nil {
		return fmt.Errorf("could not create function client dispatcher: %w", err)
	}

	cli := http.DefaultClient
	cli.Timeout = defaultHTTPTimeout

	api := client.New(log,
		client.WithHTTPClient(cli),
		client.WithHost(flagAPIEndpoint),
	)

	db, err := sql.Open(databaseDriver, flagDBConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open SQL connection: %w", err)
	}

	store, err := postgres.NewStore(db)
	if err != nil {
		return fmt.Errorf("could not create store: %w", err)
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

	queue, err := rmqConnection.OpenQueue(flagParsingQueueName)
	if err != nil {
		return fmt.Errorf("could not open redis queue: %w", err)
	}

	err = queue.StartConsuming(flagConsumerPrefetch, flagConsumerPollDuration)
	if err != nil {
		return fmt.Errorf("could not start consuming process: %w", err)
	}

	consumer := consumer.NewParsingConsumer(log, api, dispatcher, store)
	consumerName, err := queue.AddConsumer(flagRMQTag, consumer)
	if err != nil {
		return fmt.Errorf("could not add parsing consumer: %w", err)
	}

	log.Info().Str("name", consumerName).Msg("started parsing dispatcher")

	select {
	case <-sig:
		rmqConnection.StopAllConsuming()
	case err := <-failed:
		return err
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced interruption")
	}()

	log.Info().Str("name", consumerName).Msg("stopped parsing dispatcher gracefully")

	return nil
}
