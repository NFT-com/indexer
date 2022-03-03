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

	"github.com/NFT-com/indexer/job"
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
	pflag.StringVarP(&flagRMQTag, "tag", "t", "api", "watcher-web3 producer tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "network")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis url")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database")
	pflag.StringVarP(&flagDeliveryQueueName, "delivery-queue", "q", "discovery", "queue name")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "w", "parsing", "queue name")
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

	failed := make(chan error)

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * 30

	redisConnection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return err
	}

	apiClient := client.NewClient(log, client.NewOptions(
		client.WithHTTPClient(httpClient),
		client.WithHost(flagAPIEndpoint),
	))
	messageProducer, err := producer.NewProducer(redisConnection, flagDeliveryQueueName, flagParsingQueueName)
	if err != nil {
		return err
	}

	jobWatcher := watcher.NewJobWatcher(log, apiClient, messageProducer)

	discoveryJobs := make(chan job.Discovery)
	err = apiClient.SubscribeNewDiscoveryJob(discoveryJobs)
	if err != nil {
		return err
	}

	parsingJobs := make(chan job.Parsing)
	err = apiClient.SubscribeNewParsingJob(parsingJobs)
	if err != nil {
		return err
	}

	go func() {
		log.Info().Msg("job watcher starting")

		err := jobWatcher.Watch(discoveryJobs, parsingJobs)
		if err != nil {
			failed <- err
		}

		log.Info().Msg("job watcher done")
	}()

	select {
	case <-sig:
		log.Info().Msg("job watcher stopping")
		jobWatcher.Close()
		apiClient.Close()
	case err := <-failed:
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
