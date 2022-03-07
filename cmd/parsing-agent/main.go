package main

import (
	"fmt"
	"github.com/NFT-com/indexer/queue/consumer"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
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
		flagConsumerPrefetch     int64
		flagConsumerPollDuration time.Duration
		flagParsingQueueName     string
		flagLogLevel             string
	)

	pflag.StringVarP(&flagRMQTag, "tag", "t", "api", " producer tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "network")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis url")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "q", "parsing", "queue name")
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

	redisConnection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return err
	}

	queue, err := redisConnection.OpenQueue(flagParsingQueueName)
	if err != nil {
		return err
	}

	parsingConsumer, err := consumer.NewParsingConsumer(log, nil)
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

	log.Info().Str("name", consumerName).Msg("started parsing agent")

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
