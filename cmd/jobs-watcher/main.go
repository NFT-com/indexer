package main

import (
	"fmt"
	"log"
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

const (
	defaultHTTPTimeout       = time.Second * 30
	defaultDeliveryQueueName = "discovery"
	defaultParsingQueueName  = "parsing"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
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
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database number")
	pflag.StringVar(&flagDeliveryQueueName, "delivery-queue", defaultDeliveryQueueName, "name of the queue for delivery queue")
	pflag.StringVar(&flagParsingQueueName, "parsing-queue", defaultParsingQueueName, "name of the queue for parsing queue")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	failed := make(chan error)

	cli := http.DefaultClient
	cli.Timeout = defaultHTTPTimeout

	redisConnection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return fmt.Errorf("could not open connection with redis: %w", err)
	}

	api := client.New(log,
		client.WithHTTPClient(cli),
		client.WithHost(flagAPIEndpoint),
	)
	producer, err := producer.NewProducer(redisConnection, flagDeliveryQueueName, flagParsingQueueName)
	if err != nil {
		return fmt.Errorf("could not create message producer: %w", err)
	}

	watcher := watcher.New(log, api, producer)

	discoveryJobs := make(chan jobs.Discovery)
	err = api.SubscribeNewDiscoveryJob(discoveryJobs)
	if err != nil {
		return fmt.Errorf("could not subscribe to new discovery jobs: %w", err)
	}

	parsingJobs := make(chan jobs.Parsing)
	err = api.SubscribeNewParsingJob(parsingJobs)
	if err != nil {
		return fmt.Errorf("could not subscribe to new parsing jobs: %w", err)
	}

	go func() {
		log.Info().Msg("jobs watcher starting")

		err = watcher.Watch(discoveryJobs, parsingJobs)
		if err != nil {
			failed <- fmt.Errorf("could not watch jobs: %w", err)
		}

		log.Info().Msg("jobs watcher done")
	}()

	select {
	case <-sig:
		log.Info().Msg("jobs watcher stopping")
		watcher.Close()
		api.Close()
	case err = <-failed:
		log.Error().Err(err).Msg("jobs watcher aborted")
		return err
	}

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	return nil
}
