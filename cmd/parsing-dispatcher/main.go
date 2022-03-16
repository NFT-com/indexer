package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/queue/consumer"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/service/postgres"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
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
		flagDBDriver             string
		flagDBConnectionInfo     string
		flagParsingQueueName     string
		flagLambdaURL            string
		flagLogLevel             string
		flagRMQTag               string
		flagRedisDatabase        int
		flagRedisNetwork         string
		flagRedisURL             string
		flagRegion               string
		flagTestMode             bool
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.Int64VarP(&flagConsumerPrefetch, "prefetch", "p", 5, "amount of message to prefetch in the consumer")
	pflag.DurationVarP(&flagConsumerPollDuration, "poll-duration", "i", time.Second*20, "time for each consumer poll")
	pflag.StringVar(&flagDBDriver, "driver", "postgres", "name of driver to use for database connection")
	pflag.StringVarP(&flagDBConnectionInfo, "db", "d", "", "data source name for database connection")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "q", "parsing", "name of the queue for parsing")
	pflag.StringVarP(&flagLambdaURL, "function-url", "f", "", "url for the custom lambda server on local testing")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "parsing-agent", "parsing dispatcher consumer tag")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "database of the network for redis connection")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "name of the network for redis connection")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "url of the network for redis connection")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "name of the region for the lambda")
	pflag.BoolVarP(&flagTestMode, "test", "t", false, "set dispatcher component to run in test mode")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	log = log.Level(level)

	sessionConfig := aws.Config{Region: aws.String(flagRegion)}
	if flagTestMode {
		sessionConfig.Credentials = credentials.AnonymousCredentials
	}

	lambdaConfig := &aws.Config{}
	if flagLambdaURL != "" {
		lambdaConfig.Endpoint = aws.String(flagLambdaURL)
	}

	sess := session.Must(session.NewSession(&sessionConfig))
	lambdaClient := lambda.New(sess, lambdaConfig)

	dispatcher, err := function.NewClient(lambdaClient)
	if err != nil {
		return fmt.Errorf("failed to create function client dispatcher: %w", err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * 30

	apiClient := client.NewClient(log, client.NewOptions(
		client.WithHTTPClient(httpClient),
		client.WithHost(flagAPIEndpoint),
	))

	db, err := sql.Open(flagDBDriver, flagDBConnectionInfo)
	if err != nil {
		return fmt.Errorf("failed to open SQL connection: %w", err)
	}

	postgresStore, err := postgres.NewStore(db)
	if err != nil {
		return fmt.Errorf("failed to create store: %v", err)
	}

	parsingConsumer, err := consumer.NewParsingConsumer(log, apiClient, dispatcher, postgresStore)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	failed := make(chan error)

	redisConnection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return fmt.Errorf("failed to open redis connection: %w", err)
	}

	queue, err := redisConnection.OpenQueue(flagParsingQueueName)
	if err != nil {
		return fmt.Errorf("failed to open redis queue: %w", err)
	}

	err = queue.StartConsuming(flagConsumerPrefetch, flagConsumerPollDuration)
	if err != nil {
		return fmt.Errorf("failed to start consuming process: %w", err)
	}

	consumerName, err := queue.AddConsumer(flagRMQTag, parsingConsumer)
	if err != nil {
		return fmt.Errorf("failed to add parsing consumer: %w", err)
	}

	log.Info().Str("name", consumerName).Msg("started parsing dispatcher")

	select {
	case <-sig:
		redisConnection.StopAllConsuming()
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
