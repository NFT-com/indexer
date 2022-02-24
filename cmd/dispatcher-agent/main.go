package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/consumer"
	"github.com/NFT-com/indexer/dispatch/redismq"
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
		flagRMQTag               string
		flagRedisNetwork         string
		flagRedisURL             string
		flagRedisDatabase        int
		flagConsumerQueueName    string
		flagConsumerPrefetch     int64
		flagConsumerPollDuration time.Duration
		flagLogLevel             string
	)

	pflag.StringVarP(&flagRMQTag, "tag", "t", "watcher", "watcher redismq tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "network")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis url")
	pflag.IntVarP(&flagRedisDatabase, "database", "d", 1, "redis database")
	pflag.StringVarP(&flagConsumerQueueName, "queue", "q", "discovery", "queue name")
	pflag.Int64VarP(&flagConsumerPrefetch, "prefetch", "p", 5, "consumer prefetch amount")
	pflag.DurationVarP(&flagConsumerPollDuration, "poll-duration", "i", time.Second, "consumer poll duration")
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
	connection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return err
	}

	producer, err := redismq.NewProducer(connection)
	if err != nil {
		return err
	}

	queue, err := connection.OpenQueue(flagConsumerQueueName)
	if err != nil {
		return err
	}

	discoveryConsumer, err := consumer.NewDiscoveryConsumer(producer)
	if err != nil {
		return err
	}

	err = queue.StartConsuming(flagConsumerPrefetch, flagConsumerPollDuration)
	if err != nil {
		return err
	}

	consumerName, err := queue.AddConsumer(flagRMQTag, discoveryConsumer)
	if err != nil {
		return err
	}

	log.Info().Str("name", consumerName).Msg("started dispatcher agent")

	select {
	case <-sig:
		connection.StopAllConsuming()
	case err := <-failed:
		return err
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced interruption")
	}()

	return nil
}
