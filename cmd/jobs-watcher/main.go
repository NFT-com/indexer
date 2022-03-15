package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/watcher"
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
		flagAPIEndpoint       string
		flagRMQTag            string
		flagRedisNetwork      string
		flagRedisURL          string
		flagRedisDatabase     int
		flagDeliveryQueueName string
		flagParsingQueueName  string
		flagLogLevel          string
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.StringVarP(&flagRMQTag, "tag", "t", "jobs-watcher", "jobs watcher producer tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "name of the network for redis connection")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "url of the network for redis connection")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "database of the network for redis connection")
	pflag.StringVarP(&flagDeliveryQueueName, "delivery-queue", "q", "discovery", "name of the queue for delivery queue")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "w", "parsing", "name of the queue for delivery queue")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	log = log.Level(level)

	failed := make(chan error)

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * 30

	redisConnection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return fmt.Errorf("failed to open connection with redis: %w", err)
	}

	apiClient := client.NewClient(log, client.NewOptions(
		client.WithHTTPClient(httpClient),
		client.WithHost(flagAPIEndpoint),
	))
	messageProducer, err := producer.NewProducer(redisConnection, flagDeliveryQueueName, flagParsingQueueName)
	if err != nil {
		return fmt.Errorf("failed to create message producer: %w", err)
	}

	jobWatcher := watcher.NewJobWatcher(log, apiClient, messageProducer)

	discoveryJobs := make(chan jobs.Discovery)
	err = apiClient.SubscribeNewDiscoveryJob(discoveryJobs)
	if err != nil {
		return fmt.Errorf("failed to subscriber to new discovery jobs: %w", err)
	}

	parsingJobs := make(chan jobs.Parsing)
	err = apiClient.SubscribeNewParsingJob(parsingJobs)
	if err != nil {
		return fmt.Errorf("failed to subscriber to new parsing jobs: %w", err)
	}

	go func() {
		log.Info().Msg("job watcher starting")

		err = jobWatcher.Watch(discoveryJobs, parsingJobs)
		if err != nil {
			failed <- fmt.Errorf("failed to watch jobs: %w", err)
		}

		log.Info().Msg("job watcher done")
	}()

	select {
	case <-sig:
		log.Info().Msg("job watcher stopping")
		jobWatcher.Close()
		apiClient.Close()
	case err = <-failed:
		log.Error().Err(err).Msg("job watcher aborted")
		return err
	}

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	return nil
}
