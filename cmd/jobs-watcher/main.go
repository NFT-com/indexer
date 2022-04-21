package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NFT-com/indexer/service/postgres"
	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/client"
	watcher "github.com/NFT-com/indexer/watcher/jobs"
)

const (
	databaseDriver = "postgres"

	defaultHTTPTimeout       = time.Second * 30
	defaultReadDelay         = time.Millisecond * 200
	defaultDeliveryQueueName = "discovery"
	defaultParsingQueueName  = "parsing"
	defaultAdditionQueueName = "addition"
)

func main() {
	err := run()
	if err != nil {
		// TODO: Improve this mixing logging
		// https://github.com/NFT-com/indexer/issues/32
		log.Fatal(err)
	}
}

func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagAdditionQueueName string
		flagAPIEndpoint       string
		flagDBConnectionInfo  string
		flagDatabaseReadDelay time.Duration
		flagRMQTag            string
		flagRedisNetwork      string
		flagRedisURL          string
		flagRedisDatabase     int
		flagDeliveryQueueName string
		flagParsingQueueName  string
		flagLogLevel          string
	)

	pflag.StringVar(&flagAdditionQueueName, "addition-queue", defaultAdditionQueueName, "name of the queue for addition queue")
	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.StringVarP(&flagDBConnectionInfo, "database", "d", "", "data source name for database connection")
	pflag.DurationVar(&flagDatabaseReadDelay, "read-delay", defaultReadDelay, "data read for new jobs delay")
	pflag.StringVarP(&flagRMQTag, "tag", "t", "jobs-watcher", "jobs watcher producer tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.IntVar(&flagRedisDatabase, "redis-database", 1, "redis database number")
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

	producer, err := producer.NewProducer(redisConnection, flagDeliveryQueueName, flagParsingQueueName, flagAdditionQueueName)
	if err != nil {
		return fmt.Errorf("could not create message producer: %w", err)
	}

	// Open database connection.
	db, err := sql.Open(databaseDriver, flagDBConnectionInfo)
	if err != nil {
		log.Error().Err(err).Msg("could not open SQL connection")
		return err
	}

	// Create the database store.
	store, err := postgres.NewStore(db)
	if err != nil {
		log.Error().Err(err).Msg("could not create store")
		return err
	}

	watcher := watcher.New(log, api, producer, store, flagDatabaseReadDelay)

	log.Info().Msg("jobs watcher starting")
	watcher.Watch()

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

	log.Info().Msg("jobs watcher done")

	return nil
}
