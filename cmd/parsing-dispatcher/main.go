package main

import (
	"fmt"
	"github.com/NFT-com/indexer/service/client"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/queue/consumer"
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
		flagRMQTag               string
		flagRedisNetwork         string
		flagRedisURL             string
		flagRedisDatabase        int
		flagConsumerPrefetch     int64
		flagConsumerPollDuration time.Duration
		flagParsingQueueName     string
		flagTestMode             bool
		flagLambdaURL            string
		flagRegion               string
		flagLogLevel             string
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "parsing-agent", "consumer tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "network")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis url")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "q", "parsing", "queue name")
	pflag.Int64VarP(&flagConsumerPrefetch, "prefetch", "p", 5, "consumer prefetch amount")
	pflag.DurationVarP(&flagConsumerPollDuration, "poll-duration", "i", time.Second*20, "consumer poll duration")
	pflag.BoolVarP(&flagTestMode, "test", "t", false, "test mode")
	pflag.StringVarP(&flagLambdaURL, "function-url", "f", "", "lambda url")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "aws region")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return err
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
		return err
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * 30

	apiClient := client.NewClient(log, client.NewOptions(
		client.WithHTTPClient(httpClient),
		client.WithHost(flagAPIEndpoint),
	))

	parsingConsumer, err := consumer.NewParsingConsumer(log, apiClient, dispatcher)
	if err != nil {
		return err
	}

	failed := make(chan error)

	redisConnection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return err
	}

	queue, err := redisConnection.OpenQueue(flagParsingQueueName)
	if err != nil {
		return err
	}

	err = queue.StartConsuming(flagConsumerPrefetch, flagConsumerPollDuration)
	if err != nil {
		return err
	}

	consumerName, err := queue.AddConsumer(flagRMQTag, parsingConsumer)
	if err != nil {
		return err
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

	return nil
}
